package router

import (
	"asp-payment/api-server/middleware"
	"asp-payment/common/pkg/config"
	"asp-payment/common/service"
	"asp-payment/common/service/callback"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Api(app *fiber.App) {

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code": 0,
			"msg":  "Success",
			"data": "from ip:" + c.IP(),
		})
	})

	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code": 0,
			"msg":  "Success",
			"data": fmt.Sprintf("time:%s port:%d version : %s", time.Now().Format("2006-01-02 15:04:05.000000"), c.Locals("serverPort").(int), config.Version),
		})
	})

	app.Get("/dashboard", monitor.New())

	api := app.Group("/")
	{

		//管理端接口
		admin := api.Group("admin", middleware.AdminAuth())
		{
			// 代付审核
			admin.Post("/payoutAudit", func(ctx *fiber.Ctx) error {
				return service.NewPayoutServer(ctx).PayoutAuditApi()
			})
			// 查询商户渠道账户信息
			admin.Post("/merchantChannelQuery", func(ctx *fiber.Ctx) error {
				return service.NewMerchantProjectServer(ctx).MerchantAccountChannelQuery()
			})
			// 查询cp使用的支付渠道
			admin.Post("/getMerchantProjectChannel", func(ctx *fiber.Ctx) error {
				return service.NewMerchantProjectServer(ctx).GetMerchantProjectCurrentChannel()
			})
		}

		// 添加路由
		payment := api.Group("pay")
		{
			payment.Post("/order", func(ctx *fiber.Ctx) error {
				return service.NewPayOrderServer(ctx).Pay()
			})

			//收银台页面
			payment.Get("/checkout/:sn/:redirectType?/:payType?", func(ctx *fiber.Ctx) error {
				return service.NewPayOrderServer(ctx).CheckOut()
			})

			//收银台页面 QrCode
			payment.Get("/qrcode/:sn", func(ctx *fiber.Ctx) error {
				return service.NewPayOrderServer(ctx).QrCode()
			})

			// 代收查询
			payment.Get("/queryOrder", func(ctx *fiber.Ctx) error {
				return service.NewPayOrderServer(ctx).OrderQueryApi()
			})

			// 代付申请
			payment.Post("/payout", func(ctx *fiber.Ctx) error {
				return service.NewPayoutServer(ctx).Payout()
			})

			// 待收查询
			payment.Get("/queryPayout", func(ctx *fiber.Ctx) error {
				return service.NewPayoutServer(ctx).PayoutQuery()
			})

			// 查询cp账户信息
			payment.Post("/merchantQuery", func(ctx *fiber.Ctx) error {
				return service.NewMerchantProjectServer(ctx).MerchantAccountQuery()
			})

			// 添加收益人
			//payment.Post("/beneficiary", func(ctx *fiber.Ctx) error {
			//	return service.NewPayoutServer(ctx).AddBeneficiary()
			//})
		}

		// firstpay 路由组
		firstPayCallBack := api.Group("callback/firstpay")
		{
			firstPayCallBack.Post("/order", func(ctx *fiber.Ctx) error {
				return callback.NewFirstPayCallBackServer(ctx).PayOrder()
			})
			firstPayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewFirstPayCallBackServer(ctx).Payout()
			})
		}

		// zPay 路由组
		zPayCallBack := api.Group("callback/zpay")
		{
			zPayCallBack.Get("/order", func(ctx *fiber.Ctx) error {
				return callback.NewZPayCallBackServer(ctx).PayOrder()
			})
			zPayCallBack.Get("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewZPayCallBackServer(ctx).Payout()
			})
		}

		// seveneightPay 路由组
		seveneightPayCallBack := api.Group("callback/seveneight")
		{
			seveneightPayCallBack.Post("/order", func(ctx *fiber.Ctx) error {
				return callback.NewSevenEightPayCallBackService(ctx).PayOrder()
			})
			seveneightPayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewSevenEightPayCallBackService(ctx).Payout()
			})
		}

		// amarquickpay 路由组
		amarquickPayCallBack := api.Group("callback/amarquickpay")
		{
			amarquickPayCallBack.Post("/order", func(ctx *fiber.Ctx) error {
				return callback.NewAmarquickPayCallBackService(ctx).PayOrder()
			})
			amarquickPayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewAmarquickPayCallBackService(ctx).Payout()
			})
		}

		// abcPay 路由组
		abcPayCallBack := api.Group("callback/abcpay")
		{
			abcPayCallBack.Post("/order", func(ctx *fiber.Ctx) error {
				return callback.NewAbcPayCallBackService(ctx).PayOrder()
			})
			abcPayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewAbcPayCallBackService(ctx).Payout()
			})
		}

		// fynzonPay 路由组
		fynzonPayCallBack := api.Group("callback/fynzonpay")
		{
			fynzonPayCallBack.Post("/order", func(ctx *fiber.Ctx) error {
				return callback.NewFynzonPayCallBackService(ctx).PayOrder()
			})
			fynzonPayCallBack.Post("/beneficiary", func(ctx *fiber.Ctx) error {
				return callback.NewFynzonPayCallBackService(ctx).Beneficiary()
			})
			fynzonPayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewFynzonPayCallBackService(ctx).Payout()
			})
		}

		// mypay 回调路由组
		mypayCallBack := api.Group("callback/mypay")
		{
			mypayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewMyPayCallBackService(ctx).PayOrder()
			})
		}

		// haodapay 路由组
		haodaPayCallBack := api.Group("callback/haodapay")
		{
			haodaPayCallBack.Post("/order", func(ctx *fiber.Ctx) error {
				return callback.NewHaoDaPayCallBackService(ctx).PayOrder()
			})
			haodaPayCallBack.Post("/payout", func(ctx *fiber.Ctx) error {
				return callback.NewHaoDaPayCallBackService(ctx).Payout()
			})
		}

		views := api.Group("notify")
		{
			views.Get("/order/return_success", func(ctx *fiber.Ctx) error {
				return service.NewPayOrderServer(ctx).ReturnSuccess()
			})
			views.Get("/order/return_error", func(ctx *fiber.Ctx) error {
				return ctx.Render("order/return_error", fiber.Map{
					"title": "payment error",
				})
			})
		}
	}

}
