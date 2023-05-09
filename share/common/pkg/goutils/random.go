package goutils

import (
	"context"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"math/rand"
	"share/common/pkg/goRedis"
	"strings"
	"time"
)

// 随机生成字符串
func RandomString(l int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	bytes := []byte(str)
	var result []byte = make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return BytesToString(result)
}

// 随机生成纯字符串
func RandomPureString(l int) string {
	str := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	bytes := []byte(str)
	var result []byte = make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return BytesToString(result)
}

// 随机生成数字字符串
func RandomNumber(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return BytesToString(result)
}

func GenerateSerialNumBer(region string) string {
	if region == "" {
		region = "India"
	}
	preFixMap := map[string]string{
		"India":  "10", // 默认
		"hk":     "11",
		"Payout": "12", // 付款
	}
	_, ok := preFixMap[region]
	var prefix = ""
	if ok {
		prefix = preFixMap[region]
	}
	nowTime := carbon.Now().ToDateTimeString()                     // 2020-08-05 13:14:15
	nowTimeFormat := carbon.Parse(nowTime).ToShortDateTimeString() // 20200805131415
	curTime := carbon.Now().Timestamp()
	redisKey := cast.ToString(curTime)
	redisKeyPrefix := redisKey
	autoIncrement := goRedis.Redis.Incr(context.Background(), redisKeyPrefix).Val()
	// 设置过期时间
	if autoIncrement <= 1 {
		goRedis.Redis.Expire(context.Background(), redisKeyPrefix, time.Second*3)
	}
	// 获取统一毫秒内递增值
	roundNumber := cast.ToString(autoIncrement)

	var res strings.Builder
	res.WriteString(prefix)
	res.WriteString(nowTimeFormat)
	// 总的7位递增值 补0
	roundLen := 8 - len(roundNumber)
	for roundLen > 0 {
		roundLen--
		res.WriteString("0")
	}
	res.WriteString(roundNumber)

	return res.String()
}

func GenerateSerialNumBerBack(region string) string {
	if region == "" {
		region = "India"
	}
	preFixMap := map[string]string{
		"India":  "10", // 默认
		"hk":     "11",
		"Payout": "12", // 付款
	}
	_, ok := preFixMap[region]
	var prefix = ""
	if ok {
		prefix = preFixMap[region]
	}

	time := GetNowTimesTamp()
	rundNumber := RandomNumber(7)

	timeIntNano := GetDateTimeUnixNano()

	timeStringNano := Int642String(timeIntNano)

	str3 := timeStringNano[10:len(timeStringNano)]

	var res strings.Builder
	res.WriteString(prefix)
	res.WriteString(time)
	res.WriteString(str3)
	res.WriteString(rundNumber)

	return res.String()
}
