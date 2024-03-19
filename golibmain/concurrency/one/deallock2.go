package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1 // 开启子goroutine写入数据
	}()

	fmt.Println(<-ch) // 阻塞住，一旦ch有数据，则读取成功
}
