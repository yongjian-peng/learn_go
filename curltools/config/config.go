package config

import (
	"curltools/goRedis"
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	AppFileName = "app"
)

type conf struct {
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
	Server struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
		Port int    `yaml:"port"`
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
	//初始化redis
	goRedis.Init(AppConfig.Redis.Host, AppConfig.Redis.Port, AppConfig.Redis.DB, AppConfig.Redis.Auth, AppConfig.Server.Name, AppConfig.Server.Env)

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
