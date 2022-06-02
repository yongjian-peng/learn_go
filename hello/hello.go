package main

import (
	"fmt"
	"time"

	// "log"

	// "example.com/greetings"
	"example.com/morestrings"
)

func main() {
	timeUnix := time.Now().Unix() //已知的时间戳

	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	fmt.Println(formatTimeStr) //打印结果：2017-04-11 13:30:3

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timeStr)

	// log.SetPrefix("greetings:")
	// log.SetFlags(0)

	// names := []string{"Gladys", "Samantha", "Darrin"}

	// message, err := greetings.Hellos(names)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // 消息到控制台
	// fmt.Println(message)
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
}
