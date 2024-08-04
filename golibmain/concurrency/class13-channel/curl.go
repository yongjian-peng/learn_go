package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	total = 109
	limit = 20
)

func main() {
	requests := (total + limit - 1) / limit

	var wg sync.WaitGroup
	wg.Add(requests)

	for i := 0; i < requests; i++ {
		offset := i * limit
		go func(offset int) {
			defer wg.Done()

			// 构造API请求
			requestURL := constructRequestURL(offset, limit)

			fmt.Println("requestURL=>", requestURL)

			//// 发送API请求
			//response, err := sendRequest(requestURL)
			//if err != nil {
			//	fmt.Printf("请求失败：%v\n", err)
			//	return
			//}
			//
			//// 处理API响应
			//processResponse(response)
		}(offset)
	}

	wg.Wait()
}

func constructRequestURL(offset, limit int) string {
	// 根据offset和limit构造请求URL
	requestURL := fmt.Sprintf("https://example.com/api?offset=%d&limit=%d", offset+1, limit)
	return requestURL
}

func sendRequest(requestURL string) ([]byte, error) {
	// 发送API请求并返回响应
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func processResponse(response []byte) {
	// 处理API响应
	// 这里可以根据实际需求来解析和处理响应数据
	// 这里只是一个示例，简单地打印响应内容
	fmt.Println(string(response))
}
