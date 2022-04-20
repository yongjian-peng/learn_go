package main

import (
	"curltools/curl"
	"curltools/model"
	"curltools/signature"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
)

const GroupNum int = 1

const ResultChanNum int = 10
const ExitChanNum int = 10

var ResultSuccessNum int64
var wg sync.WaitGroup
var logger *log.Logger

func main() {

	for i := 0; i < GroupNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			usechannel()

		}(i)
	}
	wg.Wait()
	fmt.Printf("所有的请求已完成,请求数量：%d\n", ResultChanNum*GroupNum)
	// fmt.Printf("%x", &TotalAmount)
	// fmt.Printf("所有的请求已完成,请求数量：%d\n", ResultSuccessNum)
	resNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(ResultSuccessNum)/float64(ResultChanNum*GroupNum)), 64)
	resNum = resNum * 100
	fmt.Println("所有的请求已完成,成功率：", resNum, "%")
}

func usechannel() {
	resultChan := make(chan int, ResultChanNum)
	exitChan := make(chan bool, ExitChanNum)

	for i := 0; i < ExitChanNum; i++ {
		go calc(resultChan, exitChan)
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

func calc(resChan chan int, exitChan chan bool) {

	SendPost()

	fmt.Println("exit")
	exitChan <- true
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
		return
	}

	var postWithJosn = &model.PostWithJson{
		Appid: postjson["Appid"],
		Sn:    postjson["Sn"],
		Sign:  postJson,
	}

	url := "http://test.co"

	result := curl.Curl(url, *postWithJosn)
	// fmt.Println(result["code"])

	if result["code"] == 200.0 {
		atomic.AddInt64(&ResultSuccessNum, 1)
	}

	// Output: *goz.Response
}
