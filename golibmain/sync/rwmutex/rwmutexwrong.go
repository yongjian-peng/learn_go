package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var lock sync.RWMutex
	var wg sync.WaitGroup

	wg.Add(1)

	lock.RLock()
	fmt.Println("rlock 1...")
	wg.Done()
	// go func() {
	// 	lock.Lock()
	// 	fmt.Println("lock.....")
	// 	time.Sleep(time.Second)
	// 	lock.Unlock()
	// 	fmt.Println("unlock...")
	// 	wg.Done()
	// }()

	time.Sleep(time.Second)

	lock.RLock()
	fmt.Println("rlock2 ...")
	lock.RUnlock()
	fmt.Println("runlock 1...")
	lock.RUnlock()
	fmt.Println("runlock 2...")

	wg.Wait()
}
