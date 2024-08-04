package main

import (
	"fmt"
	"sync"
)

// 使用两个goroutine交替打印1-100之间的奇数和偶数
// 两个goroutine分别打印奇数和偶数，在两个goroutine中每次for循环时fmt都能得到执行，所以为了实现交替打印1-100，需要使用两对无缓冲的channel。
// 20行和27行构成第一对同步channel，21和29行构成第二对同步channel，
// 第一对同步channel完成读取数据之前，奇数就已经打印；第二对同步channel完成打印之前，偶数就会打印，所以最终输出就是交替打印1-100。

var ch = make(chan struct{})
var wg sync.WaitGroup

func go1() {
	defer wg.Done()
	for i := 1; i <= 5; i += 2 {
		fmt.Println("线程1", i)
		ch <- struct{}{} // 不能与上一行交换位置
		<-ch
	}
}
func go2() {
	defer wg.Done()
	for i := 2; i <= 5; i += 2 {
		<-ch
		fmt.Println("线程2", i)
		ch <- struct{}{}
	}
}

func main() {
	wg.Add(2)
	go go1()
	go go2()
	wg.Wait()
}
