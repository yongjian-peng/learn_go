package main

import (
	"curltools/curl"
	"curltools/model"
	"curltools/signature"
)

func main() {

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
