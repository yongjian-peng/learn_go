package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// Go channel 如何控制goroutine并发执行顺序
func main() {
	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)
	ch3 := make(chan struct{}, 1)
	ch1 <- struct{}{}
	wg.Add(3)
	go print("goroutine1", ch1, ch2)
	go print("goroutine2", ch2, ch3)
	go print("goroutine3", ch3, ch1)

	wg.Wait()
}

func print(gorouine string, inputchan chan struct{}, outchan chan struct{}) {
	// 模拟内部操作耗时
	time.Sleep(1 * time.Millisecond)
	select {
	case <-inputchan:
		fmt.Println("gorouine:", gorouine)
		outchan <- struct{}{}
	}
	wg.Done()
}
