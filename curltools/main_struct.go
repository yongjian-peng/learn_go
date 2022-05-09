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

const GroupNum int = 2

const ResultChanNum int = 2

var ResultSuccessNum int64
var wg sync.WaitGroup
var logger *log.Logger

type RunSend struct {
	execNum int
	ch      chan struct{}
}

func (r *RunSend) ConcurrentExec(task func(r *RunSend)) {
	for i := 0; i < r.execNum; i++ {
		go task(r)
	}

	for i := 0; i < r.execNum; i++ {
		<-r.ch
	}
	fmt.Printf("执行了【%d】个请求\n", r.execNum)
}

func main() {

	wg.Add(GroupNum)

	p := &sync.Pool{
		New: createNewClient,
	}

	for i := 0; i < GroupNum; i++ {
		go func(index int) {
			defer wg.Done()
			run := RunSend{
				execNum: ResultChanNum,
				ch:      make(chan struct{}, ResultChanNum),
			}
			run.ConcurrentExec(func(r *RunSend) {
				SendPost(p)
				r.ch <- struct{}{}
			})
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
