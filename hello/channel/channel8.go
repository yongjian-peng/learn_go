package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

/* 出自：https://go101.org/article/channel-closing.html */
/* N发送者情况的变种：必须关闭数据通道，告诉接收者数据发送结束 */
/* 在上述N-sender 情况的解决方案中，为了保持通道关闭原则，我们必须避免关闭数据通道。但是，有时需要
最后关闭数据通道，让接收方知道数据发送结束。对于这种情况，我们可以通过使用中间通道将N个
发送者的情况转换为一个发送者的情况，中间通道只有一个发送者，这样我们就可以关闭它而不是关闭原来的
数据通道
*/

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 1000000
	const NumReceivers = 10
	const NumSenders = 1000
	const NumThirdParties = 15

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	middleCh := make(chan int)
	closing := make(chan string)
	closed := make(chan struct{})

	var stoppedBy string

	stop := func(by string) {
		select {
		case closing <- by:
			<-closed
		case <-closed:
		}
	}

	go func() {
		exit := func(v int, needSend bool) {
			close(closed)
			if needSend {
				dataCh <- v
			}
			close(dataCh)
		}

		for {
			select {
			case stoppedBy = <-closing:
				exit(0, false)
				return
			case v := <-middleCh:
				select {
				case stoppedBy = <-closing:
					exit(v, true)
					return
				case dataCh <- v:
				}
			}
		}
	}()

	for i := 0; i < NumThirdParties; i++ {
		go func(id string) {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			stop("3rd-party#" + id)
		}(strconv.Itoa(i))
	}

	// 发送者
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					stop("sender#" + id)
					return
				}

				select {
				case <-closed:
					return
				default:
				}

				select {
				case <-closed:
					return
				case middleCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 接收者
	for range [NumReceivers]struct{}{} {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
