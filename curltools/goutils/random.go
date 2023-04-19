package goutils

import (
	"context"
	"curltools/config"
	"curltools/goRedis"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
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

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
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
	str := "1234567890"
	bytes := []byte(str)
	var result []byte

	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(100))))
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
	return GetSequenceId(prefix)
	//nowTime := carbon.Now().ToDateTimeString()                     // 2020-08-05 13:14:15
	//nowTimeFormat := carbon.Parse(nowTime).ToShortDateTimeString() // 20200805131415
	//curTime := carbon.Now().TimestampNano()
	//redisKey := cast.ToString(curTime)
	//redisKeyPrefix := constant.SnPrefix + redisKey
	//
	//autoIncrement := Sequence(redisKeyPrefix)
	////autoIncrement := Redis().Incr(context.Background(), redisKeyPrefix).Val()
	////// 设置过期时间
	////if autoIncrement <= 1 {
	////	Redis().Expire(context.Background(), redisKeyPrefix, time.Second*60)
	////}
	//// 获取统一毫秒内递增值
	//roundNumber := cast.ToString(autoIncrement)
	//// 截取毫秒时间戳 后九位数组
	//strNano := redisKey[10:len(redisKey)]
	//
	//var res strings.Builder
	//res.WriteString(prefix)
	//res.WriteString(nowTimeFormat)
	//res.WriteString(strNano)
	//// 总的7位递增值 补0
	//roundLen := 7 - len(roundNumber)
	//for roundLen > 0 {
	//	roundLen--
	//	res.WriteString("0")
	//}
	//res.WriteString(roundNumber)
	//
	//return res.String()
}

func GetSequenceId(prefix string) string {
	//每纳秒可生成 9 个
	//通过redis获取单位纳秒内的自增id
	nanoSecond := carbon.Now().TimestampNano()
	key := GetSequenceKey(nanoSecond)
	sequence := Sequence(key, 0)
	for sequence > 999 {
		//暂停1纳秒
		time.Sleep(1 * time.Nanosecond)
		nanoSecond = carbon.Now().TimestampNano()
		key = GetSequenceKey(nanoSecond)
		sequence = Sequence(key, 0)
	}
	if sequence == 0 {
		return ""
	}
	strNanosecond := cast.ToString(nanoSecond)
	return fmt.Sprintf("%s%s%s%s", prefix, carbon.CreateFromTimestampNano(nanoSecond).Format("ymdHis"), strNanosecond[10:17], fmt.Sprintf("%03d", sequence))
}

func GetSequenceKey(currentTime int64) string {
	return fmt.Sprintf("%s:%s:seq:%d", config.AppConfig.Server.Name, config.AppConfig.Server.Env, currentTime)
}

func Sequence(key string, num int) int64 {
	if num > 2 {
		return 0
	}
	num++
	pipe := goRedis.Redis.TxPipeline()

	//autoIncrement := goRedis.Redis.Incr(context.Background(), key).Val()
	//
	//// 设置过期时间
	//ttl := goRedis.Redis.Expire(context.Background(), key, 1*time.Second)
	//res, _ := ttl.Result()
	//if res == false {
	//	fmt.Println("resBool: ", res)
	//}
	// 执行事务操作，可以通过pipe读写redis
	_ = pipe.Incr(context.Background(), key).Val()
	pipe.Expire(context.Background(), key, 1*time.Second)

	// 通过Exec函数提交redis事务
	res1, _ := pipe.Exec(context.Background())

	//for key1, val := range res1 {
	//	item := strings.Split(val.String(), " ")
	//
	//	fmt.Println("key1:", key1, "val: ", val, "item:", item[2])
	//}
	item := strings.Split(res1[0].String(), " ")

	if len(item) < 3 {
		Sequence(key, num)
	}

	res := cast.ToInt64(item[2])
	//autoIncrement = res1[0]

	//fmt.Println("res1: ", res1)
	//fmt.Println("res1: ", res1[2])
	//fmt.Println("err1: ", err1)

	//if autoIncrement <= 1 {
	//	goRedis.Redis.Expire(context.Background(), key, 2*time.Second)
	//}
	return res
}

func GenerateSerialNumBer2(region string) string {
	if region == "" {
		region = "India"
	}
	//preFixMap := map[string]string{
	//	"India":  "10", // 默认
	//	"hk":     "11",
	//	"Payout": "12", // 付款
	//}
	//_, ok := preFixMap[region]
	//var prefix = ""
	//if ok {
	//	prefix = preFixMap[region]
	//}

	//time := GetNowTimesTamp()
	//rundNumber := RandomNumber(7)
	rundNumber := RandStr(13)

	timeIntNano := GetDateTimeNotUnixNano()
	fmt.Println("timeIntNano: ", timeIntNano)
	timeStringNano := Int642String(timeIntNano)

	//str3 := timeStringNano[10:len(timeStringNano)]

	var res strings.Builder
	//res.WriteString(prefix)
	//preTime := time[7:len(time)]
	//res.WriteString(preTime)
	res.WriteString(timeStringNano)
	res.WriteString(rundNumber)

	return res.String()
}

func Redis() *redis.Client {
	return goRedis.Redis
}
