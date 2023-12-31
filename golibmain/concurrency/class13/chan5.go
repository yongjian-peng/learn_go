package main

import (
	"fmt"
)

// 核心思想，首先发送数字小的，打印数字，然后再发送数字加一，并环形打印
func main() {
	signChan1 := make(chan struct{})
	signChan2 := make(chan struct{})
	signChan3 := make(chan struct{})
	signChan4 := make(chan struct{})

	mainSignChan := make(chan struct{})
	for i := 1; i <= 4; i++ {
		go func(i int) {
			for {
				select {
				case <-signChan1:
					fmt.Println(1)
					signChan2 <- struct{}{}
				case <-signChan2:
					fmt.Println(2)
					signChan3 <- struct{}{}
				case <-signChan3:
					fmt.Println(3)
					signChan4 <- struct{}{}
				case <-signChan4:
					fmt.Println(4)
					signChan1 <- struct{}{}
				}
			}
		}(i)

	}
	signChan1 <- struct{}{}
	<-mainSignChan
}
