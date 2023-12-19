package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	go dosomething(100, &wg) // 启动第一个goroutine
	go dosomething(110, &wg) // 启动第二个goroutine
	go dosomething(120, &wg) // 启动第三个goroutine
	go dosomething(130, &wg) // 启动第四个goroutine

	wg.Wait() // 主goroutine等待完成
	fmt.Println("Done")
}

func dosomething(millisece time.Duration, wg *sync.WaitGroup) {
	duration := millisece * time.Millisecond
	time.Sleep(duration) // 故意sleep一段时间

	fmt.Println("后台执行， duration:", duration)
	wg.Done()
}
