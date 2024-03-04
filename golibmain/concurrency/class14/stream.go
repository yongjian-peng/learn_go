package main

import (
	"fmt"
	"sync"
)

// 把channel当做流式管道使用，也就是把Channel看作流（Stream)
// 提供跳过几个元素，或者是只取其中的几个元素等方法
func asStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	s := make(chan interface{}) // 创建一个unbrffered的Channel
	go func() {                 // 启动一个goroutine,往s中塞数据
		defer close(s)             // 退出时关闭chan
		for _, v := range values { // 遍历数组
			select {
			case <-done:
				return
			case s <- v: // 将数组元素塞入到chan中
			}
		}
	}()
	return s
}

// 只取前n个流数据
func takeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{}) // 创建输出流
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream: // 从输入流中读取元素
			}
		}
	}()
	return takeStream
}

func main() {
	// 创建一个表示结束的通道
	done := make(chan struct{})
	// 创建一个 WaitGroup 来等待 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(1)
	// 调用 asStream 函数，传入 done 和一些值
	stream := asStream(done, "hello", 42, true, "haha", "hehe")

	// 使用 for range 循环从 stream 中读取值，直到通道关闭
	for value := range stream {
		fmt.Println(value) // 打印从 stream 中读取的值
	}
	wg.Done()

	// 在一段时间后发送信号到 done 通道，以结束 goroutine 和 stream
	wg.Wait()
	close(done)

	// 注意：由于 main 函数在关闭 done 后立即结束，
	// 如果没有足够的时间让 goroutine 执行完成，可能不会打印出所有的值。
	// 为了确保所有的值都被打印，我们可以添加一个同步机制，例如 WaitGroup。
}
