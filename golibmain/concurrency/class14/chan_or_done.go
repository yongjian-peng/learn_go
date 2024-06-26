package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: // 也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: // 超过两个，二分递归处理
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}

		}
	}()

	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})

	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// 任务编排
// 比如，你发送同一个请求到多个微服务节点，只要任意一个微服务节点返回结果，就算成功，这时候，可以参考
func main() {
	start := time.Now()

	<-or(
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Second),
	)

	fmt.Println("done after %v", time.Since(start))
}
