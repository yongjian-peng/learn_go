package curl

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"

	"curltools/model"
	"curltools/mylogger"

	"github.com/idoubi/goz"
)

func Curl(url string, postWithJosn *model.PostWithJson, pool *sync.Pool) (res map[string]interface{}) {

	cli := pool.Get().(*goz.Request)

	defer pool.Put(cli)

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
	fmt.Println("type:", reflect.TypeOf(body))
	fmt.Println(body)
	mylogger.WriteLog("curl.txt", string(body))
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}

	// var responseData = model.ResponseData
	// err = json.Unmarshal([]byte(body), responseData)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if responseData.Code == 200 {
	// 	fmt.Println(responseData.Code)
	// }

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err == nil {
		// fmt.Println("==============json str转map=================")
		// fmt.Println(dat["code"])
		// if dat["code"] == 200 {

		// }
	}
	res = dat
	return res
}

func CurlNew(url string, postWithJosn map[string]interface{}, Headers map[string]interface{}, pool *sync.Pool) (res map[string]interface{}) {

	cli := pool.Get().(*goz.Request)

	defer pool.Put(cli)

	resp, err := cli.Post(url, goz.Options{
		Headers: Headers,
		JSON:    postWithJosn,
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	//fmt.Println("type:", reflect.TypeOf(body))
	fmt.Println("body: ", body)
	mylogger.WriteLog("curl.txt", string(body))
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}

	// var responseData = model.ResponseData
	// err = json.Unmarshal([]byte(body), responseData)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if responseData.Code == 200 {
	// 	fmt.Println(responseData.Code)
	// }

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err == nil {
		// fmt.Println("==============json str转map=================")
		// fmt.Println(dat["code"])
		// if dat["code"] == 200 {

		// }
	}
	res = dat
	return res
}

func CurlNewGet(url string, postWithJosn map[string]interface{}, Headers map[string]interface{}, pool *sync.Pool) (res map[string]interface{}) {

	cli := pool.Get().(*goz.Request)

	defer pool.Put(cli)

	resp, err := cli.Get(url, goz.Options{
		Headers: Headers,
		JSON:    postWithJosn,
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	//fmt.Println("type:", reflect.TypeOf(body))
	fmt.Println("body: ", body)
	mylogger.WriteLog("curl.txt", string(body))
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}

	// var responseData = model.ResponseData
	// err = json.Unmarshal([]byte(body), responseData)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if responseData.Code == 200 {
	// 	fmt.Println(responseData.Code)
	// }

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err == nil {
		// fmt.Println("==============json str转map=================")
		// fmt.Println(dat["code"])
		// if dat["code"] == 200 {

		// }
	}
	res = dat
	return res
}
