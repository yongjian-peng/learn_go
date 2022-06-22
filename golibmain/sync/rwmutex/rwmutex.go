package main

import (
	"fmt"
	"sync"
)

func main() {
	var lock sync.RWMutex

	lock.Lock()
	fmt.Println("lock ...")
	lock.Unlock()

	lock.RLock()
	fmt.Println("Rlock 1")

	lock.RLock()
	fmt.Println("Rlock 2")

	lock.RUnlock()
	fmt.Println("Runlock 2")

	lock.RUnlock()
	fmt.Println("RUnlock 1")
}
