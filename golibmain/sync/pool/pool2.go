package main

import (
	"fmt"
	"sync"
)

type Job struct {
	ID int
}

type Worker struct {
	ID         int
	JobChannel chan Job
	Quit       chan bool
}

type Pool struct {
	WorkerNum   int
	JobChannel  chan Job
	WorkerQueue chan chan Job
	Quit        chan bool
	wg          sync.WaitGroup
}

// NewWorker 创建一个新的工作者
func NewWorker(id int, workerQueue chan chan Job) Worker {
	return Worker{
		ID:         id,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
	}
}

// Start 启动工作者
func (w Worker) Start(workerQueue chan chan Job) {
	go func() {
		for {
			workerQueue <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				fmt.Printf("Worker %d started job %d\n", w.ID, job.ID)
				// 执行任务
				fmt.Printf("Worker %d finished job %d\n", w.ID, job.ID)
			case <-w.Quit:
				return
			}
		}
	}()
}

// Start 启动工作者池
func (p *Pool) Start() {
	for i := 0; i < p.WorkerNum; i++ {
		worker := NewWorker(i, p.WorkerQueue)
		worker.Start(p.WorkerQueue)
	}

	go func() {
		for {
			select {
			case job := <-p.JobChannel:
				worker := <-p.WorkerQueue
				worker <- job
			case <-p.Quit:
				for i := 0; i < p.WorkerNum; i++ {
					worker := <-p.WorkerQueue
					worker <- Job{} // 发送空任务，通知协程退出
				}
				p.wg.Done()
				return
			}
		}
	}()
}

// AddJob 添加作业到作业通道
func (p *Pool) AddJob(job Job) {
	p.JobChannel <- job
}

// Stop 停止工作者池
func (p *Pool) Stop() {
	p.Quit <- true
	p.wg.Wait()
}

// 手写一个 携程池
// 引用地址：https://xie.infoq.cn/article/e8ac1b7599a0046e4d94b2454
func main() {
	pool := Pool{
		WorkerNum:   5,
		JobChannel:  make(chan Job),
		WorkerQueue: make(chan chan Job, 5),
		Quit:        make(chan bool),
	}

	pool.Start()

	for i := 0; i < 10; i++ {
		job := Job{ID: i}
		pool.AddJob(job)
	}

	pool.Stop()
}
