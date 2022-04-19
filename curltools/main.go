package main

import (
	"curltools/curl"
	"curltools/model"
	"curltools/signature"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

const GroupNum int = 4

const IntChanNum int = 500
const ResultChanNum int = 500
const ExitChanNum int = 2

var ResultSuccessNum int64
var wg sync.WaitGroup

func main() {

	for i := 0; i < GroupNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer usechannel()

		}(i)
	}
	wg.Wait()
	fmt.Printf("所有的请求已完成,请求数量：%d\n", IntChanNum*GroupNum)
	// fmt.Printf("%x", &TotalAmount)
	// fmt.Printf("所有的请求已完成,请求数量：%d\n", ResultSuccessNum)
	resNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(ResultSuccessNum)/float64(IntChanNum*GroupNum)), 64)
	resNum = resNum * 100
	fmt.Println("所有的请求已完成,成功率：", resNum, "%")
}

func usechannel() {
	intChan := make(chan int, IntChanNum)
	resultChan := make(chan int, ResultChanNum)
	exitChan := make(chan bool, ExitChanNum)

	go func() {
		for i := 0; i < IntChanNum; i++ {
			intChan <- i
		}

		close(intChan)
	}()

	for i := 0; i < ExitChanNum; i++ {
		go calc(intChan, resultChan, exitChan)
	}

	//等待所有计算的goroutine全部退出
	go func() {
		for i := 0; i < ExitChanNum; i++ {
			<-exitChan
			fmt.Println("wait goroute ", i, " exited")
		}
		close(resultChan)
	}()

	for v := range resultChan {
		fmt.Printf("result:")
		fmt.Println(v)
	}

}

func calc(taskChan chan int, resChan chan int, exitChan chan bool) {
	for v := range taskChan {

		go SendPost()

		resChan <- v
	}

	fmt.Println("exit")
	exitChan <- true
}

func addResultNum(i int) {
	atomic.AddInt64(&ResultSuccessNum, 1)
	fmt.Printf("result:%d,\n", i)
}

func SendPost() {
	postjson := map[string]string{
		"appid": "1000258000",
		"sn":    "11202204021728004314877380235",
	}

	signal := "NuOpbWyKYV9scy9CclISR1kHLANqyVlL"

	// var postWithJosn = &model.PostWithJson{
	// 	Appid: "1000258",
	// 	Sn:    "11202204021728004314877380235",
	// }

	err, postJson := signature.Signature(postjson, signal)

	if err != nil {
		return res
	}

	var postWithJosn = &model.PostWithJson{
		Appid: "100025800",
		Sn:    "11202204021728004314877380235",
		Sign:  postJson,
	}

	url := "http://api-yong.hk.blueoceantech.co/order/query"

	result := curl.Curl(url, *postWithJosn)
	// fmt.Println(result["code"])
	if result["code"] == 200.0 {
		atomic.AddInt64(&ResultSuccessNum, 1)
	}

	// Output: *goz.Response
}
