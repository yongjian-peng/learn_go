package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

/* 优雅关闭通道的解决方案 */
/* M个接收者，一个发送者，发送者通过关闭数据通道说“不再发送” */
/* 这是最简单的情况，就是让发送方在不想发送更多的时候关闭数据通道 */
// 出自：https://go101.org/article/channel-closing.html
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 100

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)

	go func() {
		for {
			if value := rand.Intn(Max); value == 0 {
				close(dataCh)
				return
			} else {
				dataCh <- value
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
