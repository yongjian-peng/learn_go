package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1
		close(ch)
	}()

	for {
		select {
		case i := <-ch:
			fmt.Println(i)
		default:
			break
		}
	}
}
