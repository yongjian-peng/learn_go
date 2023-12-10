package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const (
	mutexLocked = 1 << iota
	mutexWoken
	mutexWaiterShift = iota
)

func TestByte(t *testing.T) {
	fmt.Println("mutexLocked=>", mutexLocked)
	fmt.Println("mutexWoken=>", mutexWoken)
	fmt.Println("mutexWaiterShift=>", mutexWaiterShift)
	a := 53 // 00110101
	fmt.Printf("a = %b\n", a)
	b := 32 // 00000110
	fmt.Printf("b = %b\n", b)
	fmt.Printf("a &^ b = %b\n", a&^b) // 输出 a &^ b = 110011(二进制)
	//a &= ^b
	fmt.Printf("a = %b\n", a) // 输出 a = 110011（二进制）

	fmt.Printf("a = %d\n", a)
	fmt.Printf("b = %d\n", b)

	c := 5
	d := 0
	e := c | d
	fmt.Printf("e = %b\n", e)
	fmt.Printf("e = %d\n", e)
}

func TestCond(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)

			// 加锁更改等待条件
			c.L.Lock()
			ready++
			c.L.Unlock()

			log.Printf("运动员#%d 已准备就绪\n", i)
			// 广播唤醒所有的等待者
			c.Broadcast()
		}(i)
	}

	c.L.Lock()
	//for ready != 10 {
	c.Wait()
	log.Println("裁判员被唤醒一次")
	//}
	c.L.Unlock()

	//所有的运动员是否就绪
	log.Println("所有运动员都准备就绪。比赛开始，3，2，1, ......")
}
