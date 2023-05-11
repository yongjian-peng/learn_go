package main

import "fmt"

// 出自：https://go101.org/article/channel-closing.html
func sendData(ch chan int, done chan bool) {
	for i := 0; i < 10; i++ {
		done <- false
		ch <- i
		done <- true
	}
	// close(ch)
	// done <- true
}

func main() {
	ch := make(chan int)
	done := make(chan bool)
	go sendData(ch, done)
	for {
		select {
		case val, ok := <-ch:
			if ok {
				fmt.Println(val)
			} else {
				fmt.Println("channel closed")
			}
		case <-done:
			fmt.Println("ch ")
			close(ch)
			return
		}
	}
}
