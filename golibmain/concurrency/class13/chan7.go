package main

import (
	"fmt"
	"time"
)

func main() {
	chArr := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}

	for k, _ := range chArr {
		if k == len(chArr)-1 {
			go goon(chArr[k], chArr[0], k+1)
		} else {
			go goon(chArr[k], chArr[k+1], k+1)
		}
	}

	chArr[0] <- struct{}{}
	select {}
}

func goon(ch chan struct{}, ch2 chan struct{}, index int) {
	time.Sleep(time.Duration(index*10) * time.Millisecond)
	for {
		<-ch
		fmt.Printf("I am No %d Goroutine\n", index)
		time.Sleep(time.Second)
		ch2 <- struct{}{}
	}
}
