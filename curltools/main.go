package main

import (
	"curltools/config"
	"curltools/curl"
	"curltools/goutils"
	"fmt"
	"github.com/spf13/cast"
	"log"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/idoubi/goz"
)

const GroupNum int = 300

const ResultChanNum int = 2
const ExitChanNum int = 2

var (
	ResultSuccessNum int64
	WG2              sync.WaitGroup
	MU               *sync.Mutex
	logger           *log.Logger
	resCodes         sync.Map
)

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

	//params := GetParams()
	params := GetPayoutParams()

	signal := "hqPgg8BP06Sibmp3Q4KW0hMroUyHoOOr"
	//signal := "BessXQtS8xy3SWudR0i6gc0zfSk3Z8ll"
	//signal := "P1GN6nHEenrLFq1Qi8iAItDTulfEK05m"

	sign, err := goutils.GetSign(params, signal)
	if err != nil {
		return
	}
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	headers["appid"] = "1000012"
	//headers["appid"] = "1000026"
	//headers["appid"] = "1000028"
	headers["version"] = "1.0"
	headers["signature"] = sign

	//url := "http://sunny-pay.co/pay/order"
	//url := "http://sunny-pay.co/pay/payout"

	//url := "https://hopepatti.fun/pay/order"
	url := "https://hopepatti.fun/pay/payout"

	result := curl.CurlNew(url, params, headers, pool)
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

func GetParams() map[string]interface{} {
	//mu.Lock()
	mount := goutils.RandomNumber(5)

	amount := cast.ToInt(mount) * 100

	if amount <= 10000 {
		amount = 100000
	}

	orderId := goutils.GenerateSerialNumBer("")
	userId := goutils.RandStr(10)

	fmt.Printf("GetParams orderId：%s , userId: %s, amount %d \n", orderId, userId, amount)

	params := make(map[string]interface{})
	params["payment_method"] = "sunny.h5"
	params["order_id"] = orderId
	params["order_currency"] = "INR"
	params["order_amount"] = amount
	params["order_name"] = goutils.RandomString(10)
	params["user_id"] = userId
	//params["return_url"] = "https://11.fun/callback/zpay/order"
	//params["notify_url"] = "https://11.fun/callback/zpay/order"
	params["customer_name"] = "ericluzhonghua"
	params["customer_phone"] = "9036830689"
	params["customer_email"] = "ericluzhonghua@gmail.com"
	params["device_info"] = "device_info"
	params["order_note"] = "order_note"
	//mu.Unlock()
	return params
}

func GetPayoutParams() map[string]interface{} {
	//mu.Lock()
	mount := goutils.RandomNumber(5)

	amount := cast.ToInt(mount) * 100

	if amount <= 10000 {
		amount = 100000
	}

	orderId := goutils.GenerateSerialNumBer("Payout")
	userId := goutils.RandStr(10)

	fmt.Printf("GetParams orderId：%s , userId: %s, amount %d \n", orderId, userId, amount)

	params := make(map[string]interface{})
	params["order_id"] = orderId
	params["order_currency"] = "INR"
	params["order_amount"] = amount
	params["order_name"] = goutils.RandomString(10)
	params["user_id"] = userId
	//params["notify_url"] = "https://11.fun/callback/zpay/order"
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

	//mu.Unlock()
	return params
}
