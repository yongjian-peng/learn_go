package main

import (
	"sync"
	"time"
)

// 引用：https://stackoverflow.com/questions/36857167/how-to-correctly-use-sync-cond

func main() {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	m.Lock()
	go func() {
		m.Lock()
		time.Sleep(1 * time.Second)
		m.Unlock()
		c.Broadcast()
	}()

	time.Sleep(2 * time.Second)
	c.Wait()
	m.Unlock()
}
