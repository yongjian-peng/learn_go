package main

import (
	"context"
	"fmt"
	"log"
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

type CrondJob struct {
	Name string
	Job  func()
	Ctx  context.Context
	Wg   *sync.WaitGroup
}

type UserCrondJob struct {
}

func NewSimpleUserCrondJob() *UserCrondJob {
	return &UserCrondJob{}
}

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())

	go listenSignal(cancel)

	jobs := make([]*CrondJob, 0)

	userCrondJob := NewSimpleUserCrondJob()

	jobs = append(jobs, &CrondJob{
		Name: "syncUserJob", Job: userCrondJob.Config, Ctx: ctx, Wg: &wg,
	})

	for _, job := range jobs {
		job.Run()
	}
}

func listenSignal(cancel context.CancelFunc) {
	fmt.Println("context.Signal")
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case s := <-c:

			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("收到信号需要平滑关闭")
				// logger.Println("收到信号需要平滑关闭")
				cancel()
			default:
				cancel()
			}
		}
	}
}

func (c *CrondJob) Run() {
	fmt.Println(fmt.Sprintf("%s crondJob is Run", c.Name))
	c.Wg.Add(1)

	go func() {
		for {
			select {
			case <-c.Ctx.Done():
				// 收到信号，减少一个 goroutine 等待
				c.Wg.Done()
				return
			default:
				// 执行任务
				c.Job()
				time.Sleep(time.Second * 1)
			}
		}
	}()
}

func (j *UserCrondJob) Config() {
	fmt.Println("config")
}
