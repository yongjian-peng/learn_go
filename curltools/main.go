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

	"github.com/idoubi/goz"
)

const GroupNum int = 10

const ResultChanNum int = 10
const ExitChanNum int = 10

var ResultSuccessNum int64
var wg sync.WaitGroup
var logger *log.Logger

func main() {

	wg.Add(GroupNum)

	p := &sync.Pool{
		New: createNewClient,
	}

	for i := 0; i < GroupNum; i++ {
		go func(index int) {
			defer wg.Done()
			usechannel(p)

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

	fmt.Println("exit")
	exitChan <- true
}

func SendPost(pool *sync.Pool) {
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

	result := curl.Curl(url, *postWithJosn, pool)
	// fmt.Println(result["code"])

	if result["code"] == 200.0 {
		atomic.AddInt64(&ResultSuccessNum, 1)
	}

	// Output: *goz.Response
}
