package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

/* M个接收者，N个发送者，其中任何一个通过通知主持人关闭一个额外的信号通道来表示"让我们结束发送" */
/* 这是最复杂的情况，我们不能让任何接收者和发送者关闭数据通道，而且我们不能让任何一个接受者
关闭一个额外的信号通道来通知所有发送者和接收者退出游戏，做任何一个都会破坏通道关闭原则。
但是，我们可以引入支持人角色来关闭额外的信号通道。以下示例中的一个技巧是如何使用try-send
操作通知主持人关闭附加信号通道
*/
// 出自：https://go101.org/article/channel-closing.html
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	stopCh := make(chan struct{})

	toStop := make(chan string, 1)

	var stoppedBy string

	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
