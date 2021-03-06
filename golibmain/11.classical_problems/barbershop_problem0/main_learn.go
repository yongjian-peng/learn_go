package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	seatsLock        sync.Mutex
	seats            int
	cond             = sync.NewCond(&seatsLock)
	ResultSuccessNum int64
)

func main() {
	go barber()
	go customers()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Printf("总共来了 %d位顾客", ResultSuccessNum)
}

func randomPause(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max)))
}

// 理发师
func barber() {
	for {
		// 等待一个客户
		log.Println("Tony老师尝试请求一个顾客")
		seatsLock.Lock()
		for seats == 0 {
			cond.Wait()
		}
		seats--
		seatsLock.Unlock()
		log.Println("Tony老师找到一位顾客，开始理发")
		randomPause(2000)
	}
}

// 模拟顾客陆陆续续的过来
func customers() {
	for {
		randomPause(1000)
		go customer()
	}
}

// 顾客
func customer() {
	seatsLock.Lock()
	defer seatsLock.Unlock()

	if seats == 3 {
		log.Println("没有空闲座位了，一位顾客开始离开了")
		return
	}
	seats++
	atomic.AddInt64(&ResultSuccessNum, 1)
	cond.Broadcast()

	log.Println("一位顾客开始坐下排队理发")
}
