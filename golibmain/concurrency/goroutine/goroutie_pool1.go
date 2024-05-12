package main

import (
	"fmt"
	"sync"
)

// Task 任务结构体
type Task struct {
	ID int
	// 任务
	Job func()
}

// Pool 协程池结构体
type Pool struct {
	// 任务通道
	taskChan chan Task
	// 工作协程数量
	workerCont int
	// 等待组
	wg sync.WaitGroup
}

// NewPool 创建协程池
func NewPool(workerCount int) *Pool {
	workChannel := make(chan Task, workerCount)
	return &Pool{
		taskChan:   workChannel,
		workerCont: workerCount,
		wg:         sync.WaitGroup{},
	}
}

// SubmitTask 向携程池提交任务
func (p *Pool) SubmitTask(task Task) {
	p.taskChan <- task
	p.wg.Add(1)
}

// StartWorkers 启动工作协程
func (p *Pool) StartWorkers() {
	for i := 0; i < p.workerCont; i++ {
		go p.worker()
	}
}

// worker 工作协程
func (p *Pool) worker() {
	for task := range p.taskChan {
		defer p.wg.Done()
		fmt.Printf("Worker received task %d\n", task.ID)
		task.Job()
		fmt.Printf("Worker completed task %d\n", task.ID)
	}
}

func main() {
	// 创建一个协程池，设置工作协程数量为5
	pool := NewPool(5)

	// 提交任务到协程池
	for i := 1; i < 5; i++ {
		task := Task{
			ID: i,
			Job: func() {
				fmt.Printf("Task %d is running\n", i)
			},
		}
		pool.SubmitTask(task)
	}

	// 启动工作协程
	pool.StartWorkers()

	// 等待所有任务完成
	pool.wg.Wait()
}
