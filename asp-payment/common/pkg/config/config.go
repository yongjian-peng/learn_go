package config

import (
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	AppFileName = "app"
)

type conf struct {
	Log struct {
		StorageLocation string `yaml:"storage_location"`
		MaxAge          int    `yaml:"max_age"`
		MaxBackups      int    `yaml:"max_backups"`
		MaxSize         int    `yaml:"max_size"`
		LogLevel        string `yaml:"log_level"`
	}
	Db struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		DbName       string `yaml:"db_name"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	}
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Auth string `yaml:"auth"`
		DB   int    `yaml:"db"`
	}
	JobConfig struct {
		QueryBeforeDay int `yaml:"query_before_day"`
	} `yaml:"jobConfig"`
	Server struct {
		Name        string `yaml:"name"`
		Env         string `yaml:"env"`
		Port        int    `yaml:"port"`
		AdminSecret string `yaml:"admin_secret"`
	}
	Urls struct {
		ZpayPayoutNotifyUrl           string `yaml:"zpay_payout_notify_url"`           // zpay_payout_notify
		ZpayH5NotifyUrl               string `yaml:"zpay_h5_notify_url"`               // zpay_h5_notify
		SevenEightH5NotifyUrl         string `yaml:"seveneight_h5_notify_url"`         // seveneight_h5_notify
		SevenEightPayoutNotifyUrl     string `yaml:"seveneight_payout_notify_url"`     // seveneight_h5_notify
		AmarquickpayWappayNotifyUrl   string `yaml:"amarquickpay_wappay_notify_url"`   // amarquickpay wappay notify url
		SevenEightOrderReturnUrl      string `yaml:"seveneight_order_return_url"`      // 代收成功跳转url
		WapPayCheckoutPreUrl          string `yaml:"wap_pay_checkout_pre_url"`         // 收银台url前缀地址
		WapPayCheckoutQrCodeUrl       string `yaml:"wap_pay_checkout_qrcode_url"`      // 收银台 qrcode 前缀地址
		FynzonpayWappayNotifyUrl      string `yaml:"fynzonpay_wappay_notify_url"`      // fynzonpay 代收回调地址
		FynzonpayWappayErrorUrl       string `yaml:"fynzonpay_wappay_error_url"`       // fynzonpay 代收回调地址 失败
		FynzonpayPayoutNotifyUrl      string `yaml:"fynzonpay_payout_notify_url"`      // fynzonpay 代付回调地址
		FynzonpayBeneficiaryNotifyUrl string `yaml:"fynzonpay_beneficiary_notify_url"` // fynzonpay 受益人回调地址
		FynzonpayPayoutErrorUrl       string `yaml:"fynzonpay_payout_error_url"`       // fynzonpay 代付回调地址 失败
		FynzonpayOrderReturnUrl       string `yaml:"fynzonpay_order_return_url"`       // fynzonpay 成功跳转地址
		AbcPayH5NotifyUrl           string `yaml:"abcPay_h5_notify_url"`           // abcPay_h5_notify
		AbcPayPayoutNotifyUrl       string `yaml:"abcPay_payout_notify_url"`       // abcPay_payout_notify
		AbcPayOrderReturnUrl        string `yaml:"abcpay_order_return_url"`        // 代收成功跳转url
	}
}

var (
	AppConfig = conf{}
)

const Version = "V1.1.2"

// InitConfig api config
func InitConfig() {
	//获取配置文件
	initConfig()
}

func initConfig() {
	//获取配置文件
	filePath := flag.String("c", "./config.yaml", "asp-payment config path")
	port := flag.Int("port", 0, "asp-payment port")
	flag.Parse()
	data, err := os.ReadFile(*filePath)
	if err != nil {
		panic("Conf file ReadFile:" + err.Error())
	}
	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		panic("Conf file ReadFile:" + err.Error())
	}
	if *port != 0 {
		AppConfig.Server.Port = *port
	}
	//初始化数据库
	database.Init(AppConfig.Db.User, AppConfig.Db.Password, AppConfig.Db.Host, AppConfig.Db.Port, AppConfig.Db.DbName, AppConfig.Db.MaxIdleConns, AppConfig.Db.MaxOpenConns)
	//初始化redis
	goRedis.Init(AppConfig.Redis.Host, AppConfig.Redis.Port, AppConfig.Redis.DB, AppConfig.Redis.Auth, AppConfig.Server.Name, AppConfig.Server.Env)
	//初始化数据仓库
	repository.Init()
	//设置log配置
	logger.SetLogConfig(logger.LogConfig{
		StorageLocation: AppConfig.Log.StorageLocation,
		MaxAge:          AppConfig.Log.MaxAge,
		MaxBackups:      AppConfig.Log.MaxBackups,
		MaxSize:         AppConfig.Log.MaxSize,
		LogLevel:        AppConfig.Log.LogLevel,
	})
}

func IsDevEnv() bool {
	if AppConfig.Server.Env == "Dev" {
		return true
	}
	return false
}

func IsTestEnv() bool {
	if AppConfig.Server.Env == "Test" {
		return true
	}
	return false
}
