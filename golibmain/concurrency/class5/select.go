package main

import (
	"fmt"
	"sync"
)

func foo(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("in foo")
}

func goo(wg *sync.WaitGroup) {
	defer wg.Done()
	var i int
	fmt.Println("in goo", i)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go goo(&wg)
	go foo(&wg)

	wg.Wait()

	//time.Sleep(time.Second)
	//select {}
}
