package curl

import (
	"fmt"
	"log"

	"curltools/model"

	"github.com/idoubi/goz"
)

func Curl(url string, postWithJosn model.PostWithJson) {
	// var postWithJosn = model.PostWithJson{
	// 	Appid: "Hao Chen",
	// 	Sn:    "Male",
	// 	Sign:  "abc",
	// }
	// postWithJosn.Appid = "1000258"
	// postWithJosn.Sn = "11202204021728004314877380235"
	// postWithJosn.Sign = "abc"

	cli := goz.NewClient()

	resp, err := cli.Post(url, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: postWithJosn,
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Println(body)
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}
}
