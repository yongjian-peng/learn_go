package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	type key int
	quit := false
	ctx, cancel := context.WithCancel(context.TODO())
	k := key(100)
	ctx = context.WithValue(ctx, k, 100)
	go func() {
		time.Sleep(time.Second)
		cancel() // parent context被cancel后，其子context也会被取消
	}()

	for {
		if quit {
			break
		}
		select {
		case <-ctx.Done(): // 子context被取消，结束循环
			quit = true
		default:
			fmt.Println("value=>", ctx.Value(k))
			time.Sleep(time.Second)
		}
	}

	for {
		for i := 0; i < 10; i++ {
			if i == 8 {
				break
			}
			fmt.Println("i=>", i)
		}
		//return
		break
	}

	fmt.Println("haha")
}
