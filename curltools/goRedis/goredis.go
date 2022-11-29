package goRedis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Ctx = context.Background()
var env string
var serviceName string

// 解耦
// 高内聚 低耦合
func Init(Host string, Port int, Db int, Password string, appName, appEnv string) {
	env = appEnv
	serviceName = appName
	masterAddress := fmt.Sprintf("%s:%d", Host, Port)
	Redis = redis.NewClient(&redis.Options{
		Addr:     masterAddress,
		Password: Password, // no password set
		DB:       Db,       // use default DB
	})

	// 检测心跳
	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatalln("main connect redis failed")
	}
}

func Close() {
	_ = Redis.Close()
}

func GetKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", serviceName, env, key)
}

func Lock(lockName string) bool {
	return Redis.SetNX(context.Background(), GetKey(lockName), 1, 30*time.Second).Val()
}

func UnLock(lockName string) {
	Redis.Del(context.Background(), GetKey(lockName))
}
