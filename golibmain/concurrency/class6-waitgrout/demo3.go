package main

import (
	"fmt"
	"time"
)

// 不带缓存的，不加goroutine 先写后读，会报错，如果先读后写，不加入goroutine, 会报错
// 带缓存的， 先写后读，不会报错，

func main() {
	//noCacheBeginWriteThenReader() // deadlock
	//noCacheBeginReaderThenWrite() // deadlock
	//CacheBeginWriteThenReader() // OK
	CacheBeginReaderThenWrite() // OK 但是，goroutine 需要处理执行完成，可以
}

func noCacheBeginWriteThenReader() {
	c := make(chan string)
	c <- "aa"

	s, ok := <-c
	if ok {
		fmt.Println("ok")
	}

	fmt.Println("s=>", s)

	fmt.Println("chan")
}

func noCacheBeginReaderThenWrite() {
	c := make(chan string, 1)

	s, ok := <-c
	if ok {
		fmt.Println("ok")
	}

	c <- "aa"

	fmt.Println("s=>", s)

	fmt.Println("chan")

}

func CacheBeginWriteThenReader() {
	c := make(chan string, 1)

	go func() {
		time.Sleep(3 * time.Second)
		c <- "aa"
	}()

	s, ok := <-c
	if ok {
		fmt.Println("ok")
	}

	fmt.Println("s=>", s)

	fmt.Println("chan")
}

func CacheBeginReaderThenWrite() {
	c := make(chan string, 1)

	go func() {
		s, ok := <-c
		if ok {
			fmt.Println("ok")
		}
		fmt.Println("s=>", s)
		time.Sleep(3 * time.Second)
	}()

	c <- "aa"

	fmt.Println("chan")
}
