package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var lock sync.RWMutex

	s := make(map[int]int)
	for i := 0; i < 100; i++ {
		go func(i int) {
			lock.Lock()
			s[i] = i
			lock.Unlock()
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			lock.RLock()
			fmt.Printf("map第%d 个元素是%d", i, s[i])
			lock.RUnlock()
		}(i)
	}

	time.Sleep(1 * time.Second)
}
