package main

import (
	"fmt"
	"sync"
)

type Token struct {
}

// 开启两个goroutine按照顺序打印数字
func main() {
	chs := []chan Token{make(chan Token), make(chan Token)}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go first(chs[0], chs[1], &wg)
	go second(chs[1], chs[0], &wg)

	chs[0] <- struct{}{}
	wg.Wait()

}

func first(ch chan Token, nextCh chan Token, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 100; i += 2 {
		token := <-ch
		fmt.Println("first i is :", i)
		nextCh <- token
	}
	<-ch
}

func second(ch chan Token, nextCh chan Token, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= 100; i += 2 {
		token := <-ch
		fmt.Println("second i is :", i)
		nextCh <- token
	}
}
