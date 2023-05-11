package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// 出自：https://go101.org/article/channel-closing.html
/* M个接收者，一个发送者"情况的变种：关闭请求由第三方goroutine发出" */
/* 有时，需要关闭信号必须由第三方 goroutine 发出，对于这种情况看，我们可以使用额外的信号通道来
通知发送方关闭数据通道。
*/

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 100
	const NumThirdParties = 15

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	closing := make(chan struct{})
	closed := make(chan struct{})

	stop := func() {
		select {
		case closing <- struct{}{}:
			<-closed
		case <-closed:
		}
	}

	for i := 0; i < NumThirdParties; i++ {
		go func() {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			stop()
		}()
	}

	go func() {
		defer func() {
			close(closed)
			close(dataCh)
		}()

		for {
			select {
			case <-closing:
				return
			default:
			}

			select {
			case <-closing:
				return
			case dataCh <- rand.Intn(Max):
			}
		}
	}()

	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}
