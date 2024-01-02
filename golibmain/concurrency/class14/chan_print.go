package main

import (
	"fmt"
	"time"
)

type Token struct {
}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch // 取得令牌
		fmt.Println("cent=>", id+1)
		time.Sleep(time.Second)
		nextCh <- token
	}
}

// 上节课的作业 数据传递
// 有4个goroutine,编号为1,2,3,4. 每秒钟会有一个 goroutine 打印出它自己的编号，
// 要求编写程序，让输出的编号总是按照1,2,3,4,1,2,3,4.... 这个顺序打印出来
func main() {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}

	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	// 首先把令牌交给第一个worker
	chs[0] <- struct{}{}

	select {}
}
