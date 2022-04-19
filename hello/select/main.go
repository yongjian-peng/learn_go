package main

import (
	"fmt"
	"time"
)

func write(ch1 chan int, ch2 chan int) {
	for i := 0; i < 10; i++ {
		ch1 <- i
		time.Sleep(time.Second)
		ch2 <- i * 2
		time.Sleep(time.Second)
	}
}

func main() {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	go write(ch1, ch2)

	for {
		select {
		case v := <-ch1:
			fmt.Println(v)
		case v := <-ch2:
			fmt.Println(v)
		default:
			fmt.Println("get data timeout")
			time.Sleep(time.Second)
		}
	}
}
