package main

import (
	"fmt"
	"sync"
)

var chan1, chan2 = make(chan struct{}), make(chan struct{})

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i += 2 {
			<-chan1
			fmt.Println("one ", i)
			chan2 <- struct{}{}
		}
		<-chan1
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= 100; i += 2 {
			<-chan2
			fmt.Println("two ", i)
			chan1 <- struct{}{}
		}
	}()

	chan1 <- struct{}{}

	wg.Wait()
}
