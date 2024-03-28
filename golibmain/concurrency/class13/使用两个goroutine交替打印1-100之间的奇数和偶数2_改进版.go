package main

import (
	"fmt"
	"sync"
)

// 使用两个goroutine交替打印1-100之间的奇数和偶数
// 能够做到两个协程交替打印奇偶数，而且还能保证顺序。
// 但是要注意代码22行的“<- ch1”，如果不加的话会报错，因为不加22行的代码，go2最终会被阻塞，wg.Done()无法被执行。

var ch1, ch2 = make(chan struct{}), make(chan struct{})
var wg sync.WaitGroup

func go1() {
	defer wg.Done()
	for i := 1; i <= 6; i += 2 {
		<-ch1
		fmt.Println("线程1", i)
		ch2 <- struct{}{}
	}
	<-ch1 // 不能遗漏
}

func go2() {
	defer wg.Done()
	for i := 2; i <= 6; i += 2 {
		<-ch2
		fmt.Println("线程2", i)
		ch1 <- struct{}{}
	}

}

func main() {
	wg.Add(2)
	go go1()
	go go2()
	ch1 <- struct{}{}
	wg.Wait()
}
