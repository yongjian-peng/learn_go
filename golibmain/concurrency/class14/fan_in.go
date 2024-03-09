package main

import (
	"fmt"
	"reflect"
	"time"
)

// 扇入模式，多个输入，一个输出，反射方式
func fanInRefeflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选取一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

// 扇入模式，多个输入，一个输出，递归方式
func fanInRec(chans ...<-chan interface{}) <-chan interface{} {
	switch len(chans) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(
			fanInRec(chans[:m]...),
			fanInRec(chans[m:]...),
		)
	}
}

// 两个chan，合并成一个chan,只处理两个chan
func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil { // 只要还有可读的chan
			select {
			case v, ok := <-a:
				if !ok { // a 已关闭，设置为nil
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // b 已关闭，设置为nil
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}

func main() {

	// 创建三个输入通道
	chan1 := make(chan interface{}, 1)
	chan2 := make(chan interface{}, 1)
	chan3 := make(chan interface{}, 1)

	// 扇入模式的输出通道
	outputChan := fanInRec(chan1, chan2, chan3)

	// 向三个通道发送数据
	go func() {
		chan1 <- "Hello from chan1"
		chan1 <- "Another message from chan1"
		close(chan1)
	}()

	go func() {
		chan2 <- "Hello from chan2"
		time.Sleep(1 * time.Second) // 模拟延迟
		chan2 <- "Delayed message from chan2"
		close(chan2)
	}()

	go func() {
		chan3 <- "Hello from chan3"
		close(chan3)
	}()

	// 从输出通道读取数据并打印
	for v := range outputChan {
		fmt.Println(v)
	}

}
