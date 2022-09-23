package main

import (
	"time"
	"fmt"
	"strings"
)

func main() {
	currentTime:=time.Now() 

	fmt.Println(currentTime)

	timeUnix:=time.Now().Unix()            //单位s,打印结果:1491888244

	fmt.Println(timeUnix)
 
 	timeUnixNano:=time.Now().UnixNano()  //单位纳秒,打印结果：1491888244752784461

 	fmt.Println(timeUnixNano)

	str := "1662537039581152400"

	fmt.Println(len(str))

	str3 := str[10 : len(str)]
    fmt.Println(str3) 

	//  s = strings.Split("1662537039581152400", "|")
	//  fmt.Println(s, len(s))


	 s := strings.Split("", "")
	 fmt.Println(s, len(s))
 
	 s = strings.Split("abc,abc", "")
	 fmt.Println(s, len(s))
 
	 s = strings.Split("", ",")
	 fmt.Println(s, len(s))
 
	 s = strings.Split("abc,abc", ",")
	 fmt.Println(s, len(s))
 
	 s = strings.Split("abc,abc", "|")
	 fmt.Println(s, len(s))
 
	 fmt.Println(len(""))
	 fmt.Println(len([]string{""}))
 
	//  str := ""
	//  fmt.Println(str[0])
}