package main

import (
	"fmt"
	"math/rand"
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

	randNum := rand.Int31n(5)

	fmt.Println("randNum: ", randNum)

	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	go write(ch1, ch2)

	for {
		select {
		case v := <-ch1:
			fmt.Println(v)
			randNum := rand.Int31n(3)

			fmt.Println("randNum: ", randNum)
		case v := <-ch2:
			randNum := rand.Int31n(3)

			fmt.Println("randNum: ", randNum)
			fmt.Println(v)
		default:
			randNum := rand.Int31n(3)

			fmt.Println("randNum: ", randNum)
			fmt.Println("get data timeout")
			time.Sleep(time.Second)
		}
	}
}
