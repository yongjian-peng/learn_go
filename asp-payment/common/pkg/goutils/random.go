package goutils

import (
	"asp-payment/common/pkg/goRedis"
	"context"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"math/rand"
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

func GenerateSerialNumBer(region string, serverName string, serverEnv string) string {
	if region == "" {
		region = "India"
	}
	preFixMap := map[string]string{
		"India":     "10", // 默认
		"Benefiary": "11", // 收益人标识映射
		"Payout":    "12", // 付款
	}
	_, ok := preFixMap[region]
	var prefix = ""
	if ok {
		prefix = preFixMap[region]
	}
	return GetSequenceId(prefix, serverName, serverEnv)

	//nowTime := carbon.Now().ToDateTimeString()                     // 2020-08-05 13:14:15
	//nowTimeFormat := carbon.Parse(nowTime).ToShortDateTimeString() // 20200805131415
	//curTime := carbon.Now().Timestamp()
	//redisKey := cast.ToString(curTime)
	//redisKeyPrefix := constant.SnPrefix + redisKey
	//autoIncrement := goRedis.Redis.Incr(context.Background(), redisKeyPrefix).Val()
	//// 设置过期时间
	//if autoIncrement <= 1 {
	//	goRedis.Redis.Expire(context.Background(), redisKeyPrefix, time.Second*3)
	//}
	//// 获取统一毫秒内递增值
	//roundNumber := cast.ToString(autoIncrement)
	//
	//var res strings.Builder
	//res.WriteString(prefix)
	//res.WriteString(nowTimeFormat)
	//// 总的7位递增值 补0
	//roundLen := 8 - len(roundNumber)
	//for roundLen > 0 {
	//	roundLen--
	//	res.WriteString("0")
	//}
	//res.WriteString(roundNumber)
	//
	//return res.String()
}

// GetSequenceId 生成唯一编号 系统压力 150000 的时候会有重复的值，目前可以忽略，原因是 redis 链接数的问题。后期再优化
func GetSequenceId(prefix string, serverName string, serverEnv string) string {
	//每纳秒可生成 9 个
	//通过redis获取单位纳秒内的自增id
	nanoSecond := carbon.Now().TimestampNano()
	key := GetSequenceKey(nanoSecond, serverName, serverEnv)
	sequence := Sequence(key, 0)
	for sequence > 999 {
		//暂停1纳秒
		time.Sleep(1 * time.Nanosecond)
		nanoSecond = carbon.Now().TimestampNano()
		key = GetSequenceKey(nanoSecond, serverName, serverEnv)
		sequence = Sequence(key, 0)
	}
	if sequence == 0 {
		return ""
	}
	strNanosecond := cast.ToString(nanoSecond)
	return fmt.Sprintf("%s%s%s%s", prefix, carbon.CreateFromTimestampNano(nanoSecond).Format("ymdHis"), strNanosecond[10:17], fmt.Sprintf("%03d", sequence))
}

func GetSequenceKey(currentTime int64, serverName string, serverEnv string) string {
	return fmt.Sprintf("%s:%s:seq:%d", serverName, serverEnv, currentTime)
}

func Sequence(key string, num int) int64 {
	if num > 2 {
		return 0
	}
	num++

	pipe := goRedis.Redis.TxPipeline()
	// 执行事务操作，可以通过pipe读写redis
	_ = pipe.Incr(context.Background(), key).Val()
	pipe.Expire(context.Background(), key, 1*time.Second)
	// 通过Exec函数提交redis事务
	res1, _ := pipe.Exec(context.Background())
	item := strings.Split(res1[0].String(), " ")

	if len(item) < 3 {
		Sequence(key, num)
	}
	res := cast.ToInt64(item[2])
	//autoIncrement = res1[0]
	//fmt.Println("res1: ", res1)
	//fmt.Println("res1: ", res1[2])
	//fmt.Println("err1: ", err1)
	return res
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
