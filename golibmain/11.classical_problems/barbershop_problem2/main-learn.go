package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SemaphoreLearn chan struct{}

func (s SemaphoreLearn) AcquireLearn() {
	s <- struct{}{}
}

func (s SemaphoreLearn) TryAcquireLearn() bool {
	select {
	case s <- struct{}{}: // 还有空位置
		return true
	default: // 没有空位置，顾客离开
		return false
	}
}

func (s SemaphoreLearn) ReleaseLearn() {
	<-s
}

var stats = make(SemaphoreLearn, 10)

func main() {
	// 初始化理发师 三位
	go barberLearn("Tony")
	go barberLearn("Kevin")
	go barberLearn("Allen")

	go customersLearn()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}

func randomPauseLearn(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max)))
}

// 理发师
func barberLearn(name string) {
	for {
		// 等待一个顾客
		log.Println(name + "老师尝试请求一个顾客")
		stats.ReleaseLearn()
		log.Println(name + "老师找到一个顾客，开始理发")

		randomPauseLearn(2000)
	}
}

// 模拟顾客陆陆续续的过来
func customersLearn() {
	for {
		randomPauseLearn(1000)

		go customerLearn()
	}
}

// 顾客
func customerLearn() {
	if ok := stats.TryAcquireLearn(); ok {
		log.Println("一位顾客开始坐下排队理发")
	} else {
		log.Println("没有空闲座位了，一个顾客离开了")
	}
}
