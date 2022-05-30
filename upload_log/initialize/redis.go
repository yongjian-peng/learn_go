package initialize

// package redis

import (
	"upload_log/global"

	"github.com/go-redis/redis"
	// "go.uber.org/zap"
	"log"
)

func Redis() {
	redisCfg := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Panic("reids err", err)
		// global.GLOBAL_LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
	} else {
		log.Println("redis pong", pong)
		// zap.String("pong", pong)
		// global.GLOBAL_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.Redis = client
	}
}
