package main

import (
	"fmt"
	"time"
)

/*
有一道经典的使用channel进行任务编排的题，可以尝试做一下：
有四个goroutine，编号为1,2,3,4 每秒钟会有一个 goroutine 打印出它自己的编号，要求你编写一个程序，让输出的编号是按照1,2,3,4,1,2,3,4
的顺序打印出来。
*/
func main() {
	chArr := [4]chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}

	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				<-chArr[i%4]
				fmt.Printf("i am %d\n", i+1)

				time.Sleep(1 * time.Second)
				chArr[(i+1)%4] <- struct{}{}
			}
		}(i)
	}

	chArr[0] <- struct{}{}

	time.Sleep(10 * time.Second)
}
