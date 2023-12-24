package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// withValue1()
	withValue2()
}

func withValue1() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "0001")
	ctx = context.WithValue(ctx, "key2", "0002")
	ctx = context.WithValue(ctx, "key3", "0003")
	ctx = context.WithValue(ctx, "key4", "0004")

	fmt.Println("key1=>", ctx.Value("key1"))
}

func withValue2() {
	// 创建一个父级的context
	parContext := context.Background()

	// 创建一个带有取消功能的子级context
	ctx, cancel := context.WithCancel(parContext)

	// 启动一个goroutine，执行任务
	go perfomrTask(ctx)

	// 模拟一段时间后取消任务
	time.Sleep(2 * time.Second)
	cancel()

	// 等待任务完成
	<-ctx.Done()

	// 检查任务是否成功完成或被取消
	switch ctx.Err() {
	case context.Canceled:
		fmt.Println("任务被取消")
	case context.DeadlineExceeded:
		fmt.Println("任务超时")
	default:
		fmt.Println("任务完成")
	}

	// 创建一个带有超时功能的context
	outCtx, outCancel := context.WithTimeout(parContext, 3*time.Second)
	defer outCancel()

	// 在指定时间内执行任务
	go perfomrTask(outCtx)

	// 等待任务完成或者超时
	select {
	case <-outCtx.Done():
		switch outCtx.Err() {
		case context.Canceled:
			fmt.Println("out任务被取消")
		case context.DeadlineExceeded:
			fmt.Println("out任务超时")
		default:
			fmt.Println("out任务完成")
		}
	}
}

func perfomrTask(ctx context.Context) {
	// 模拟执行任务
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			// 检查是否接收到取消信号
			fmt.Println("接收到取消信号，停止任务")
			return
		default:
			// 执行任务
			fmt.Println("执行任务=>", i)
			time.Sleep(1 * time.Second)
		}
	}
}
