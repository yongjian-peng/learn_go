package main

import (
	"curltools/curl"
	"curltools/model"
	"curltools/signature"
	"fmt"
	"sync"
	"sync/atomic"
)

const GroupNum int = 1

const IntChanNum int = 10
const ResultChanNum int = 10
const ExitChanNum int = 8

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
	fmt.Printf("所有的请求已完成")
	// fmt.Printf("%x", &TotalAmount)
	println(ResultSuccessNum)
}

func usechannel() {
	intChan := make(chan int, IntChanNum)
	resultChan := make(chan map[string]interface{}, ResultChanNum)
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

func calc(taskChan chan int, resChan chan map[string]interface{}, exitChan chan bool) {
	for v := range taskChan {

		var result map[string]interface{}
		result = SendPost()
		// fmt.Printf("result:code:%d\n", result["code"])
		if result["code"] == 200.0 {
			addResultNum(v)
			// fmt.Println(result)
			resChan <- result
		}
	}

	fmt.Println("exit")
	exitChan <- true
}

func addResultNum(i int) {
	atomic.AddInt64(&ResultSuccessNum, 1)
	fmt.Printf("result:%d,\n", i)
}

func SendPost() (res map[string]interface{}) {
	postjson := map[string]string{
		"appid": "1000258",
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
		Appid: "1000258",
		Sn:    "11202204021728004314877380235",
		Sign:  postJson,
	}

	url := "http://api.hk.blueoceantech.co/order/query"

	result := curl.Curl(url, *postWithJosn)
	// fmt.Println(result["code"])
	// if result["code"] == 200 {
	// 	atomic.AddInt64(&ResultSuccessNum, 1)
	// }

	res = result

	return res

	// Output: *goz.Response
}