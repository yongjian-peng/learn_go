package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

type MysqlLearn struct {
	weighted int64
	*semaphore.Weighted
}

type CacheLearn struct {
	weighted int64
	*semaphore.Weighted
}

func main() {
	var wg sync.WaitGroup
	s := semaphore.NewWeighted(10)

	m := MysqlLearn{weighted: 2, Weighted: s} // mysql比较慢，服务只能并发5次
	c := CacheLearn{weighted: 1, Weighted: s} // cache 比较快，服务只能并发10次
	wg.Add(16)

	// 模拟并发请求
	var (
		successMysql int32 = 8
		successCache int32 = 8
	)

	defer func() {
		fmt.Printf("success mysql:%d, sucess cached:%d: total weighted:%d\n", successMysql, successCache, successMysql*2+successCache)
	}()

	defer wg.Wait()

	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()
			b := c.TryAcquire(c.weighted)
			if !b {
				atomic.AddInt32(&successCache, -1)
				fmt.Printf("cache acquired fail:%d\n", c.weighted)
				return
			}

			time.Sleep(time.Second)
			c.Release(c.weighted)
		}()
	}

	// 模拟并发请求
	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()

			b := m.TryAcquire(m.weighted)
			if !b {
				atomic.AddInt32(&successMysql, -1)
				fmt.Printf("mysql acquired fail:%d\n", m.weighted)
				return
			}
			time.Sleep(time.Second)

			m.Release(m.weighted)
		}()
	}
}
