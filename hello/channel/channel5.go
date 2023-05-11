package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

/* 一个接收者，N个发送者，唯一的接受者通过关闭一个额外的信号通道说“请停止发送更多” */
/* 这种情况比上述情况稍微复杂一些，我们不能让接受者关闭数据通道来停止数据传输，因为这样做会破坏通道关闭原则。但是我们可以让接受者关闭一个额外的信号通道来通知发送者停止发送值 */
// 出自：https://go101.org/article/channel-closing.html
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	dataCh := make(chan int)
	stopCh := make(chan struct{})

	for i := 0; i < NumReceivers; i++ {
		go func() {
			for {
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- rand.Intn(Max):
				}
			}
		}()
	}

	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			if value == Max-1 {
				close(stopCh)
				return
			}
			log.Println(value)
		}
	}()

	wgReceivers.Wait()
}
