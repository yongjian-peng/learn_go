package global

import (
	"log"
	"upload_log/config"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Log    *log.Logger
	Gorm   *gorm.DB
	Viper  *viper.Viper
	Config config.Server
	Redis  *redis.Client
)
