package main

import (
	"fmt"
	"time"
)

func main() {
	//deadlock1()
	//deadlock2()
	//deadlock3()
	//deadlock4()
	deadlock5()
}

func deadlock1() {
	ch := make(chan int)
	ch <- 3 // 这里会发生一直阻塞的情况，执行不到下面的一句
}

func deadlock2() {
	ch := make(chan int)
	ch <- 3 // 这里会发生一直阻塞的情况，执行不到下面的一句
	//num := <-ch
	//fmt.Println("num=>", num)

	go func() {
		num := <-ch
		fmt.Println("num=>", num)
	}()
	time.Sleep(time.Second)
}

// 缓存channel写入超过缓存区数量
func deadlock3() {
	ch := make(chan int, 3)
	ch <- 3
	ch <- 4
	ch <- 5
	//ch <- 6 // 这里会发生一直阻塞的情况
}

func deadlock4() {
	//ch := make(chan int)
	ch := make(chan int, 1)
	fmt.Println(<-ch) // 这里会发生一直阻塞
}

// 多个协程互相等待
func deadlock5() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for {
			select {
			case num := <-ch1:
				fmt.Println("num=>", num)
				ch2 <- 100
			}
		}
	}()

	for {
		select {
		case num := <-ch2:
			fmt.Println("num=>", num)
			ch1 <- 300
		}
	}
}
