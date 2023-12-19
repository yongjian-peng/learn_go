package main

import (
	"sync"
	"time"
)

func main() {
	// WaitGroup 虽然可以重用，但是要有一个前提，那就是必须等到上一轮的
	// wait完成之后，才能重用 WaitGroup 执行下一轮的 Add/wait
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond)
		wg.Done()
		wg.Add(1)
	}()
	wg.Wait() // 主goroutine等待，有可能和第14行并发执行
}
