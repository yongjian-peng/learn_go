package main

import (
	"fmt"
	"time"
)

// go 内建的函数，close cap len 都可以操作 chan 类型
// close 会把 chan 关闭掉， cap 返回 chan 的容量， len 返回 chan 中缓存的还未被取走的元素的数量
func main() {
	var ch = make(chan int, 10)
	for i := 0; i < 10; i++ {
		select {
		case ch <- i:
		case v := <-ch:
			fmt.Println("v=>", v)
		}
	}

	// 可以使用 for-range 读取 chan 中的值
	go func() {
		for v := range ch {
			fmt.Println("range v=>", v)
		}
	}()

	// 或者是忽略读取的值，只是清空 chan:
	go func() {
		for range ch {
		}
	}()
	lenCh := len(ch)
	fmt.Println("len=>", lenCh)

	time.Sleep(3 * time.Second)

}
