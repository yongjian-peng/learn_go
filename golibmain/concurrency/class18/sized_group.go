package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-pkgz/syncs"
)

// SizedGroup 内部是使用信号量和 WaitGroup 实现的，它通过信号量控制并发的 goroutine 数量，或者是不控制 goroutine 数量，只控制子任务并发执行时候的数量（通过）
func main() {
	// 设置goroutine数是10
	swg := syncs.NewSizedGroup(10)

	var c uint32

	// 执行1000个子任务，只会有10个goroutine去执行
	for i := 0; i < 1000; i++ {
		swg.Go(func(ctx context.Context) {
			time.Sleep(5 * time.Millisecond)
			atomic.AddUint32(&c, 1)
		})
	}

	// 等待任务完成
	swg.Wait()
	// 输出结果
	fmt.Println(c)
}
