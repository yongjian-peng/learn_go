package main

import (
	"fmt"
	"time"
)

// ErrorOne map对象必须在使用之前初始化，如果不初始化就赋值的话，会出现panic异常
func ErrorOne() {
	var m map[int]int
	m[100] = 100
}

func One() {
	var m map[int]int
	fmt.Println("m=>", m[100])
}

type Counter struct {
	Website      string
	Start        time.Time
	PageCounters map[string]int
}

// ErrorTwo map对象必须在使用之前初始化，如果不初始化就赋值的话，会出现panic异常
func ErrorTwo() {
	var c Counter
	c.Website = "baidu.com"
	c.PageCounters["/"]++
}

// ErrorThree map并发读写
func ErrorThree() {
	var m = make(map[int]int, 10) // 初始化一个map
	go func() {
		for {
			m[1] = 1 // 设置key
		}
	}()

	go func() {
		for {
			_ = m[2] // 访问这个map
		}
	}()

	select {}
}

func main() {
	//ErrorOne()
	//One()
	//ErrorTwo()
	ErrorThree()
}
