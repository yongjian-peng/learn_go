package main

import (
	"fmt"
	"time"
)

// process 当chan不是unbufer状态时，它就不会因为没有receiver 就被block住了。
func process(timeout time.Duration) bool {
	ch := make(chan bool)

	go func() {
		// 模拟处理耗时的业务
		time.Sleep(timeout + time.Second)
		ch <- true // block
		fmt.Println("exit goroutine")
	}()

	select {
	case result := <-ch:
		return result
	case <-time.After(timeout):
		return false

	}
}
