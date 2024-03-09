package main

import (
	"fmt"
	"sync"
)

type Token struct {
}

// 开启两个goroutine按照顺序打印数字
func main() {
	chs := []chan Token{make(chan Token, 1), make(chan Token, 1)}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go first(chs[0], chs[1], &wg)
	go second(chs[1], chs[0], &wg)

	chs[0] <- struct{}{}
	wg.Wait()
}

func first(ch chan Token, nextCh chan Token, wg *sync.WaitGroup) {

	for i := 1; i < 9999999; i += 2 {
		token := <-ch
		fmt.Println("i is :", i)
		nextCh <- token
	}
	wg.Done()
}

func second(ch chan Token, nextCh chan Token, wg *sync.WaitGroup) {

	for i := 2; i < 9999999; i += 2 {
		token := <-ch
		fmt.Println("i is :", i)
		nextCh <- token
	}
	wg.Done()
}
