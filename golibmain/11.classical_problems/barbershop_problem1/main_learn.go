package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 需求：理发师和顾客的交互
// 实现：使用不超时的 和 监听停止信号处理
// 效果：和真实的理发店里面 顾客来 理发师理发 按照时间先后顺序来
type SemaphoreLean chan struct{}

func (s SemaphoreLean) AcquireLearn() {
	s <- struct{}{}
}

// 检查是否有空位置
func (s SemaphoreLean) TryAcquireLearn() bool {
	select {
	case s <- struct{}{}: // 还有空位置
		return true
	default: // 没有空位置，客户离开
		return false
	}
}

// 往channel 里面发信号
func (s SemaphoreLean) ReleaseLearn() {
	<-s
}

var seatsLearn = make(SemaphoreLean, 3)

func main() {
	go barberLearn()

	go customersLearn()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}

func randomPauseLearn(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max)))
}

// 理发师
func barberLearn() {
	for {
		// 等待一个客户
		log.Println("Tony老师尝试请求一个顾客")
		seatsLearn.ReleaseLearn()
		log.Println("Tony老师找到一位顾客，开始理发")

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

// 顾客 来的频率和理发的频率 不对等 所以用户离开的比较多
func customerLearn() {
	if ok := seatsLearn.TryAcquireLearn(); ok {
		log.Println("一位顾客开始坐下来排队理发")
	} else {
		log.Println("没有空闲位置了，一位顾客离开了")
	}
}
