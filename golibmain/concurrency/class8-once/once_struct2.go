package main

import (
	"fmt"
	"sync"
	"time"
)

var ThreeOnce struct {
	v int
	sync.Once
}

func main() {
	for i := 3; i < 100; i++ {
		go func(item int) {
			ThreeOnce.Do(func() {
				fmt.Println("先进来的值：", item)
				ThreeOnce.v += item
			})
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("v:", ThreeOnce.v)
}
