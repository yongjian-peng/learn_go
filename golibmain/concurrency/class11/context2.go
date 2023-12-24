package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)
	child := context.WithValue(ctx, "name", "peter")
	go func() {
		for {
			select {
			case <-child.Done():
				fmt.Println("it's over")
				// 放在goroutine外判断错误类型是取不到其他的错误类型的，应该放到goroutine内部
				switch child.Err() {
				case context.Canceled:
					fmt.Println("任务被取消")
				case context.DeadlineExceeded:
					fmt.Println("任务超时")
				default:
					fmt.Println("任务完成")
				}
				return
			default:
				res := child.Value("name")
				fmt.Println("name:=>", res)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()

	time.Sleep(5 * time.Second)
}
