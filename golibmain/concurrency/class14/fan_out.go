package main

import "fmt"

func main() {
	cha1 := make(chan struct{}, 1)

	over := make(chan struct{}, 1)
	listens := []chan struct{}{
		make(chan struct{}, 1),
		make(chan struct{}, 1),
		make(chan struct{}, 1),
	}
	fanOut(cha1, listens, false)
	cha1 <- struct{}{}
	cha1 <- struct{}{}
	cha1 <- struct{}{}

	close(cha1)

	go func() {
		for i := 0; i < len(listens); i++ {
			item := i
			go func(item int) {
				for {
					_, ok := <-listens[item]
					if !ok {
						over <- struct{}{}
						close(over)
						return
					}
					fmt.Printf("第%d监听着收到信息\n", item)
				}
			}(item)
		}
	}()
	<-over
}

// 扇出模式 一个输入，多个输出
func fanOut(ch <-chan struct{}, out []chan struct{}, async bool) {
	go func() {
		defer func() { // 退出时关闭所有的输出chan
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		for v := range ch { // 从输入 chan 中读取数据
			v := v
			for i := 0; i < len(out); i++ {
				item := i
				if async {
					go func() { // 异步
						out[item] <- v // 放入到输出chan中，异步方式
					}()
				} else {
					out[item] <- v // 放入到输出chan中，同步方式
				}
			}
		}
	}()
}
