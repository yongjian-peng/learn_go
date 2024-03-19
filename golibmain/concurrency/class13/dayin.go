package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 使用两个goroutine交替打印1-100之间的奇数和偶数
// runtime.Gosched这个函数的作用是让当前goroutine让出CPU，好让其它的goroutine获得执行的机会。
// 同时，当前的goroutine也会在未来的某个时间点继续运行。
// 注意：
// 当打印的数字比较大时，无法交替打印

func main() {
	// 设置可同时使用的CPU核数为1
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i < 3010; i++ {
			// 奇数
			if i%2 == 1 {
				fmt.Println("线程1打印:", i)
			}
			// 让出cpu
			// runtime.Gosched()
			time.Sleep(time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i < 3010; i++ {
			// 偶数
			if i%2 == 0 {
				fmt.Println("线程2打印:", i)
			}
			// 让出cpu
			// runtime.Gosched()
			time.Sleep(time.Millisecond)
		}
	}()
	wg.Wait()
}
