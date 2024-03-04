package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4) // 预先设定WaitGroup的计数值

	go dosomething(100, &wg)
	go dosomething(100, &wg)
	go dosomething(100, &wg)
	go dosomething(100, &wg)

	wg.Wait() // 主goroutine等待
	fmt.Println("Done")
}

func dosomething(millisecs time.Duration, wg *sync.WaitGroup) {
	duration := millisecs * time.Millisecond

	time.Sleep(duration)

	fmt.Println("后台执行，duration:=>", duration)
	wg.Done()
}
