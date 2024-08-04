package main

import (
	"context"
	"fmt"
	"time"
)

func tt(ctx context.Context, c1, c2 *chan int) {
	for {
		select {
		case n := <-*c1:
			fmt.Println("I am n=>", n)
			nn := n + 1
			if n == 4 {
				nn = 1
			}
			*c2 <- nn
			fmt.Printf("c1:%p,c2:%p\n", c1, c2)
		case <-ctx.Done():
			return
		}
	}
}

func PrintInfo() {
	ctx, cancel := context.WithCancel(context.Background())

	c1, c2, c3, c4 := make(chan int, 2), make(chan int, 2), make(chan int, 2), make(chan int, 2)
	fmt.Printf("c1:%p,c2:%p,c3:%p,c4:%p\n", &c1, &c2, &c3, &c4)
	go tt(ctx, &c1, &c2)
	go tt(ctx, &c2, &c3)
	go tt(ctx, &c3, &c4)
	go tt(ctx, &c4, &c1)

	c1 <- 1
	fmt.Println("Hello,世界")
	time.Sleep(time.Millisecond * 70)
	cancel()
	fmt.Println("Hello,世界")
}

func main() {
	PrintInfo()
}
