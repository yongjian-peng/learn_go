package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestUnLock(t *testing.T) {
	fmt.Printf("mutexWoke=>%b, %d\n", mutexWoke, mutexWoke)
	fmt.Printf("mutexWoke=>%b, %d\n", 4|2, 4|2)
	fmt.Println("old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0")
	old := 8
	fmt.Printf("old=>%b,%d\n", old, old)
	fmt.Printf("mutexWaiterShi=>%b,%d\n", mutexWaiterShi, mutexWaiterShi)
	new := old >> mutexWaiterShi
	fmt.Printf("old >> mutexWaiterShi=>%b, %d\n", new, new)
	fmt.Printf("mutexLocke|mutexWoke=>%b,%d\n", mutexLocke|mutexWoke, mutexLocke|mutexWoke)
	fmt.Printf("old&(mutexLocke|mutexWoke)=>%b, %d\n", old&(mutexLocke|mutexWoke), old&(mutexLocke|mutexWoke))

	fmt.Println(" new = (old - 1<<mutexWaiterShift) | mutexWoken")
	fmt.Printf("1<<mutexWaiterShift=>%b, %d\n", 1<<mutexWaiterShi, 1<<mutexWaiterShi)
	fmt.Printf("old - 1<<mutexWaiterShift=>%b, %d\n", old-1<<mutexWaiterShi, old-1<<mutexWaiterShi)
	fmt.Printf("old - 1<<mutexWaiterShift|mutexWoke=>%b, %d\n", (old-1<<mutexWaiterShi)|mutexWoke, (old-1<<mutexWaiterShi)|mutexWoke)
	fmt.Printf("mutexWoke=>%b, %d\n", mutexWoke, mutexWoke)
}

type Counter struct {
	sync.Mutex
	Count int
}

func TestTwoLock(t *testing.T) {
	var c Counter
	c.Lock()
	defer c.Unlock()
	c.Count++
	foo(c) // 复制锁
}

func foo(c Counter) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("In foo")
}
