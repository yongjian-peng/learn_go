package main

import (
	"curltools/curl"
	"curltools/model"
	"curltools/signature"
	"fmt"
	"sync"
)

const NUM int = 4

func main() {
	wg := &sync.WaitGroup{}
	for i := 0; i < NUM; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer SendPost()

		}(i)
	}
	wg.Wait()
	fmt.Printf("所有的请求已完成")
}

func SendPost() {
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
		return
	}

	var postWithJosn = &model.PostWithJson{
		Appid: "1000258",
		Sn:    "11202204021728004314877380235",
		Sign:  postJson,
	}

	url := "http://api.hk.blueoceantech.co/order/query"

	curl.Curl(url, *postWithJosn)

	// Output: *goz.Response
}
