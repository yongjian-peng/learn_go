package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

// goroutine 泄漏 chan 阻塞。导致goroutine 占内存.
func main() {
	fmt.Println("before goroutines:", runtime.NumGoroutine())
	block1()
	block2()
	block3()
	block4()
	block5()
	block6()
	time.Sleep(time.Second * 1)

	wg.Wait()
	fmt.Println("after goroutines:", runtime.NumGoroutine())
}

func block1() {
	var ch chan int
	for i := 0; i < 10; i++ {
		go func() {
			<-ch
		}()
	}
}

// 发送数量超过 channel 接受数量，就会造成阻塞。
func block2() {
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			ch <- 1
		}()
	}
}

// 接受数量超过发送数量，就会造成阻塞。
func block3() {
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			<-ch
		}()
	}
}

func requestWithNoClose() {
	_, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Printf("error occurred while fetching page, errors: %s", err.Error())
	}

	// defer resp.Body.Close()
}

// block4 阻塞 请求的 http body 未关闭，导致阻塞.
func block4() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			requestWithNoClose()
		}()
	}
}

// block5 互斥锁忘记解锁.
func block5() {
	var mutex sync.Mutex
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
		}()
	}
}

// block6 sync.WaitGroup 使用不当.
func block6() {
	var wg6 sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func() {
			wg6.Add(2)
			wg6.Done()
			wg6.Wait()
		}()
	}
}
