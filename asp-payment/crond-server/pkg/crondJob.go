package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CrondJob struct {
	Name string
	Job  func()
	Ctx  context.Context
	Wg   *sync.WaitGroup
}

func (c *CrondJob) Run() {
	fmt.Println(fmt.Sprintf("%s crondJob is Run", c.Name))
	c.Wg.Add(1)
	go func() {
		for {
			select {
			case <-c.Ctx.Done():
				//收到信号，减少一个goroutine等待
				c.Wg.Done()
				return
			default:
				//执行任务
				c.Job()
				//完成任务休息一秒钟
				time.Sleep(time.Second * 1)
			}
		}
	}()
}
