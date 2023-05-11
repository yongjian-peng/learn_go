package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var logger *log.Logger

func init() {
	file := "./" + time.Now().Format("20220419") + ".txt"

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "[qSkiptool]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出

	return
}

func main() {

	// 一个发送 多个接收 使用channel 分发消息
	// 注意平滑重启中的 分发和接收
	// 先收到平滑重启的信号，然后再关闭channel
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 100

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	closedCh := make(chan bool)

	go func(closedCh chan bool) {
		for {
			select {
			case close := <-closedCh:
				if close == true {
					return
				}
			default:
				value := rand.Intn(Max)
				// str := "in" + strconv.Itoa(value)
				// logger.Println(str)
				dataCh <- value
			}
		}
	}(closedCh)
	go listenSignal(dataCh, closedCh)
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				// str := "out" + strconv.Itoa(value)
				// logger.Println(str)
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}

func listenSignal(dataCh chan int, closedCh chan bool) {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case s := <-c:

			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("收到信号需要平滑关闭")
				// logger.Println("收到信号需要平滑关闭")
				closedCh <- true
				close(dataCh)
			default:
				closedCh <- true
				close(dataCh)
			}
		}
	}
}
