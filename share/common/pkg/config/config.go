package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"os"
	"share/common/pkg/database"
	"share/common/pkg/goRedis"
	"share/common/pkg/logger"
	"share/common/repository"
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
	Server struct {
		Name      string `yaml:"name"`
		Env       string `yaml:"env"`
		Port      int    `yaml:"port"`
		AssertUrl string `yaml:"assert_url"`
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
	filePath := flag.String("c", "./config.yaml", "share config path")
	port := flag.Int("port", 0, "share port")
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
