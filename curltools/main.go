package main

import (
	"curltools/config"
	"curltools/constant"
	"curltools/curl"
	"curltools/goutils"
	"fmt"
	"github.com/spf13/cast"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/idoubi/goz"
)

const GroupNum int = 1

const ResultChanNum int = 1
const ExitChanNum int = 1

var (
	ResultSuccessNum int64
	WG2              sync.WaitGroup
	MU               *sync.Mutex
	logger           *log.Logger
	resCodes         sync.Map
)

var urls = map[string]map[string]string{
	"local": {
		"order":       "http://sunny-pay.co/pay/order",
		"payout":      "http://sunny-pay.co/pay/payout",
		"query":       "http://sunny-pay.co/pay/queryOrder",
		"queryPayout": "http://sunny-pay.co/pay/queryPayout",
		"returnUrl":   "http://127.0.0.1:8084/api/v1/signup",
		"notifyUrl":   "http://127.0.0.1:8084/api/v1/signup",
	},
	"test": {
		"order":       "https://hopepatti.fun/pay/order",
		"payout":      "https://hopepatti.fun/pay/payout",
		"query":       "https://hopepatti.fun/pay/queryOrder",
		"queryPayout": "https://hopepatti.fun/pay/queryPayout",
		"returnUrl":   "http://8.219.100.35:8084/api/v1/notify",
		"notifyUrl":   "http://8.219.100.35:8084/api/v1/notify",
	},
}

var projectList = map[string]map[string]string{
	"1000012": {
		"key":   "hqPgg8BP06Sibmp3Q4KW0hMroUyHoOOr",
		"appid": "1000012",
	},
	"1000026": {
		"key":   "BessXQtS8xy3SWudR0i6gc0zfSk3Z8ll",
		"appid": "1000026",
	},
	"1000028": {
		"key":   "P1GN6nHEenrLFq1Qi8iAItDTulfEK05m",
		"appid": "1000028",
	},
	"1000259": {
		"key":   "nAZX8exqnHrzIc9r5uq01hRGzup6pskL",
		"appid": "1000259",
	},
	"1000030": {
		"key":   "n2CAjF8XvqgiZbiaWCvspFC64PAAGWrm",
		"appid": "1000030",
	},
}

func main() {

	config.InitConfig()

	WG2.Add(GroupNum)

	p := &sync.Pool{
		New: createNewClient,
	}

	for i := 0; i < GroupNum; i++ {
		go func(index int) {
			defer WG2.Done()
			usechannel(p)
		}(i)
		//go func() {
		//	fmt.Println(goutils.RandStr(10))
		//	WG2.Done()
		//}()
	}
	WG2.Wait()
	fmt.Printf("所有的请求已完成,请求数量：%d\n", ResultChanNum*GroupNum)
	// fmt.Printf("%x", &TotalAmount)
	// fmt.Printf("所有的请求已完成,请求数量：%d\n", RESUSTSUCCESSNUM)
	resNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(ResultSuccessNum)/float64(ResultChanNum*GroupNum)), 64)
	resNum = resNum * 100
	fmt.Println("所有的请求已完成,成功率：", resNum, "%")

	resCodes.Range(func(k, v interface{}) bool {
		fmt.Println("所有的请求已完成,成功编码:", k, v)
		return true
	})
	//fmt.Println("所有的请求已完成,成功编码：", resCodes, "%")
}

func createNewClient() interface{} {
	return goz.NewClient()
}

func usechannel(pool *sync.Pool) {
	resultChan := make(chan int, ResultChanNum)
	exitChan := make(chan bool, ExitChanNum)

	for i := 0; i < ExitChanNum; i++ {
		go calc(resultChan, exitChan, pool)
	}

	//等待所有计算的goroutine全部退出
	go func() {
		for i := 0; i < ExitChanNum; i++ {
			<-exitChan
			fmt.Println("wait goroute ", i, " exited")
		}
		close(resultChan)
		close(exitChan)
	}()

	for v := range resultChan {
		fmt.Printf("result:")
		fmt.Println(v)
	}

}

func calc(resChan chan int, exitChan chan bool, pool *sync.Pool) {

	SendPost(pool)

	//fmt.Println("exit")
	exitChan <- true
	// return
}

func SendPost(pool *sync.Pool) {
	// test 1000026  1000028 1000012 1000259
	// local 1000026  1000028 1000012 1000030
	// test local
	//result, err := GetOrderParams(pool, urls["local"], projectList["1000026"])
	//result, err := GetPayoutParams(pool, urls["local"], projectList["1000026"])
	//result, err := GetOrderQuery(pool, urls["local"], projectList["1000012"])
	result, err := GetPayoutQuery(pool, urls["local"], projectList["1000012"])
	if err != nil {
		return
	}
	// fmt.Println(result["code"])

	resCode := cast.ToString(result["code"])

	ccc, ok := resCodes.Load(resCode)
	if ok {
		ss := cast.ToInt(ccc) + 1
		resCodes.Store(resCode, ss)
	} else {
		resCodes.Store(resCode, 1)
	}

	//codeType := reflect.TypeOf(resCode)
	//fmt.Println("codeType: ", codeType)

	if result["code"] == 200.0 {
		atomic.AddInt64(&ResultSuccessNum, 1)
	}
	// Output: *goz.Response
}

func GetOrderParams(pool *sync.Pool, urlInfo, projectInfo map[string]string) (map[string]interface{}, error) {
	//mu.Lock()
	mount := goutils.RandomNumber(5)

	amount := cast.ToInt(mount) * 100

	if amount <= 10000 {
		amount = 100000
	}

	orderId := goutils.GenerateSerialNumBer("")
	userId := goutils.RandStr(10)

	//fmt.Printf("GetOrderParams orderId：%s , userId: %s, amount %d \n", orderId, userId, amount)

	params := make(map[string]interface{})
	params["payment_method"] = "sunny.h5"
	params["order_id"] = orderId
	params["order_currency"] = "INR"
	params["order_amount"] = amount
	params["order_name"] = goutils.RandomString(10)
	params["user_id"] = userId
	params["return_url"] = urlInfo["returnUrl"]
	params["notify_url"] = urlInfo["notifyUrl"]
	params["customer_name"] = "ericluzhonghua"
	params["customer_phone"] = "9036830689"
	params["customer_email"] = "ericluzhonghua@gmail.com"
	params["device_info"] = "device_info"
	params["order_note"] = "order_note"

	signal := projectInfo["key"]
	//signal := "BessXQtS8xy3SWudR0i6gc0zfSk3Z8ll"
	//signal := "P1GN6nHEenrLFq1Qi8iAItDTulfEK05m"

	sign, err := goutils.GetSign(params, signal)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	headers["appid"] = projectInfo["appid"]
	headers["version"] = "1.0"
	headers["signature"] = sign
	url := urlInfo["order"]

	result := curl.CurlNew(url, params, headers, pool)
	return result, nil
}

func GetPayoutParams(pool *sync.Pool, urlInfo, projectInfo map[string]string) (map[string]interface{}, error) {
	//mu.Lock()
	mount := goutils.RandomNumber(5)

	amount := cast.ToInt(mount) * 100

	if amount <= 10000 {
		amount = 100000
	}

	orderId := goutils.GenerateSerialNumBer("Payout")
	userId := goutils.RandStr(10)

	//fmt.Printf("GetOrderParams orderId：%s , userId: %s, amount %d \n", orderId, userId, amount)

	params := make(map[string]interface{})
	params["order_id"] = orderId
	params["order_currency"] = "INR"
	params["order_amount"] = amount
	params["order_name"] = goutils.RandomString(10)
	params["user_id"] = userId
	params["notify_url"] = urlInfo["notifyUrl"]
	params["customer_name"] = "ericluxxxx"
	params["customer_phone"] = "903683xxxx"
	params["customer_email"] = "ericluxxxx@gmail.com"
	params["device_info"] = "device_info"
	params["order_note"] = "order_note"
	params["ifsc"] = "ICIC000xxxx"
	params["bank_card"] = "66310170xxxx"
	params["bank_code"] = "ICICI"
	params["vpa"] = "VPA"
	params["pay_type"] = "bank"
	params["address"] = "F2/251-A SANGAM"
	params["city"] = "MAHIPALPUR EXTENSION"
	signal := projectInfo["key"]

	sign, err := goutils.GetSign(params, signal)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	headers["appid"] = projectInfo["appid"]
	headers["version"] = "1.0"
	headers["signature"] = sign
	url := urlInfo["payout"]

	result := curl.CurlNew(url, params, headers, pool)
	//mu.Unlock()
	return result, nil
}

func GetOrderQuery(pool *sync.Pool, urlInfo, projectInfo map[string]string) (map[string]interface{}, error) {

	//orderId := goutils.GenerateSerialNumBer("Payout")

	randNu := rand.Intn(100)

	orderId := constant.SnList[randNu]

	params := make(map[string]interface{})

	signal := projectInfo["key"]

	sign, err := goutils.GetSign(params, signal)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	headers["appid"] = projectInfo["appid"]
	headers["version"] = "1.0"
	headers["signature"] = sign
	url := urlInfo["query"] + "?sn=" + orderId
	result := curl.CurlNewGet(url, params, headers, pool)
	//mu.Unlock()
	return result, nil
}

func GetPayoutQuery(pool *sync.Pool, urlInfo, projectInfo map[string]string) (map[string]interface{}, error) {

	//orderId := goutils.GenerateSerialNumBer("Payout")

	randNu := rand.Intn(100)

	orderId := constant.SnPayoutList[randNu]

	params := make(map[string]interface{})

	signal := projectInfo["key"]

	sign, err := goutils.GetSign(params, signal)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	headers["appid"] = projectInfo["appid"]
	headers["version"] = "1.0"
	headers["signature"] = sign
	url := urlInfo["queryPayout"] + "?sn=" + orderId
	result := curl.CurlNewGet(url, params, headers, pool)
	//mu.Unlock()
	return result, nil
}
