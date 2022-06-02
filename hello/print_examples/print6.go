package main

import (
	"fmt"
	"sync"
)

// 3个函数分别打印 cat dog fish 要求每一个函数要起一个goroutine
// 按照 cat dog fish 顺序打印在屏幕上10次
func main() {
	var wg sync.WaitGroup
	dogCh := make(chan struct{}, 1)
	defer close(dogCh)
	catCh := make(chan struct{}, 1)
	defer close(catCh)
	fishCh := make(chan struct{}, 1)
	defer close(fishCh)

	wg.Add(3)

	go catPrint(&wg, &catCh, &dogCh)
	go dogPrint(&wg, &dogCh, &fishCh)
	go fishPrint(&wg, &fishCh, &catCh)
	catCh <- struct{}{}
	wg.Wait()
}

func catPrint(wg *sync.WaitGroup, catCh *chan struct{}, dogCh *chan struct{}) {
	count := 0
	for {
		if count >= 10 {
			wg.Done()
			fmt.Println("cat quit")
			return
		}
		<-*catCh
		fmt.Println("cat", count)
		count++
		*dogCh <- struct{}{}
	}
}

func dogPrint(wg *sync.WaitGroup, dogCh *chan struct{}, fishCh *chan struct{}) {
	count := 0
	for {
		if count >= 10 {
			wg.Done()
			fmt.Println("dog quit")
			return
		}
		<-*dogCh
		fmt.Println("dog", count)
		count++
		*fishCh <- struct{}{}
	}
}

func fishPrint(wg *sync.WaitGroup, dogCh *chan struct{}, catCh *chan struct{}) {
	count := 0
	for {
		if count >= 10 {
			wg.Done()
			fmt.Println("fish quit")
			return
		}
		<-*dogCh
		fmt.Println("fish", count)
		count++
		*catCh <- struct{}{}
	}
}
