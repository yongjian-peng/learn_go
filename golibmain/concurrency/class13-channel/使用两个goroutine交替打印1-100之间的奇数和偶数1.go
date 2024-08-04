package main

import (
	"fmt"
	"sync"
)

// 使用两个goroutine交替打印1-100之间的奇数和偶数
// 首先因为变量ch是一个无缓冲的channel, 所以只有读写同时就绪时才不会阻塞。
// 所以两个goroutine会同时进入各自的 if 语句（此时 i 是相同的），但是此时只能有一个 if 是成立的，
// 不管哪个goroutine快，都会由于读channel或写channel导致阻塞，因此程序会交替打印1-100且有顺序。

func main() {

	ch := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i < 101; i++ {
			ch <- struct{}{}
			// 奇数
			if i%2 == 1 {
				fmt.Println("线程1打印:", i)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i < 101; i++ {
			<-ch
			// 偶数
			if i%2 == 0 {
				fmt.Println("线程2打印:", i)
			}
		}
	}()

	wg.Wait()
}
