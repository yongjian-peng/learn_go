package main

import (
	"asp-payment/api-server/middleware"
	"asp-payment/api-server/router"
	"asp-payment/common/pkg/config"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/jet"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// 提交验证发布版本 更新common 关联 0419-001

	//初始化配置
	config.InitConfig()

	//模板引擎
	engine := jet.New("./views", ".jet.html").Delims("[[", "]]").Layout("body")

	if config.IsTestEnv() || config.IsDevEnv() {
		// Disable this in production
		engine.Reload(true)
	}

	//创建fiber 添加测试
	app := fiber.New(fiber.Config{
		AppName:      "API SERVER",
		ServerHeader: "sunny_api_server",
		ErrorHandler: middleware.ErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		//EnableTrustedProxyCheck: true, // 这里是属于代理 后期需要负载均衡了，再次打开配置
		TrustedProxies: []string{"0.0.0.0"},
		ProxyHeader:    fiber.HeaderXForwardedFor,
		JSONEncoder:    sonic.Marshal,
		JSONDecoder:    sonic.Unmarshal,
		Views:          engine,         //模板引擎
		ViewsLayout:    "layouts/main", //模板布局
	})

	// Setup static files
	app.Static("/assets", "./views/assets")
	//使用中间件
	app.Use(favicon.New())
	app.Use(requestid.New())
	app.Use(compress.New())
	app.Use(middleware.ServerPort())
	app.Use(middleware.Skip())
	app.Use(cors.New(middleware.Cors()))
	app.Use(logger.New(middleware.Logger()))
	app.Use(recover.New(middleware.Recover(config.IsDevEnv())))

	//api路由映射
	router.Api(app)
	//启动服务
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", config.AppConfig.Server.Port)); err != nil {
			log.Panic(err)
		}
	}()
	//平滑关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	fmt.Printf("time:%s port:%d version : %s  ApiServer shutting down... \n", time.Now().Format("2006-01-02 15:04:05.000000"), config.AppConfig.Server.Port, config.Version)
	_ = app.Shutdown()
	fmt.Printf("time:%s port:%d version : %s  ApiServer shutdown successful \n", time.Now().Format("2006-01-02 15:04:05.000000"), config.AppConfig.Server.Port, config.Version)

}
