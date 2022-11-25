package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type info struct {
	result int
}

func work() (int, error) {
	return 13, nil
}

func main() {
	var data info

	var err error
	data.result, err = work() //ok
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("info: %+v\n", data) //prints: info: {result:13}

	sn := "12456465"

	nsn := sn[0:2]
	fmt.Println(nsn)

	var params = make(map[string]interface{}, 2)
	params["aas"] = 1
	params["as"] = "as"

	for k, v := range params {
		fmt.Println("k--", k)
		fmt.Println("v--", v)
	}

	tmp := make(map[string]string)
	for k, v := range params {
		val, err := ToString(v)
		if err != nil {
			fmt.Println(err)
		}
		tmp[k] = val
	}

	fmt.Println("tmp--------", tmp)

	//获取时间戳
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	tm := time.Unix(timestamp, 0)
	fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))
	fmt.Println(tm.Format("02/01/2006 15:04:05 PM"))

}

func ToString(val interface{}) (str string, err error) {
	// fmt.Printf("%T", val)
	var s string
	if vv, ok := val.(float64); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(float32); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(int); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(int64); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(string); ok {
		s = string(vv)
	} else {
		return s, errors.New("不支持的参数类型")
	}
	return s, nil
}
