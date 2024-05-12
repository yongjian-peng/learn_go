package main

import (
	"context"
	"sync"
)

// Semaphore 使用waitGroup和channel实现并发同时控制最大并发数
// 参考golang.org/x/sync.errgroup实现返回err功能
type Semaphore struct {
	c       chan struct{}
	wg      sync.WaitGroup
	cancel  func()
	errOnce sync.Once
	err     error
}

func NewSemaphore(maxSize int) *Semaphore {
	return &Semaphore{
		c: make(chan struct{}, maxSize),
	}
}

func NewSemaphoreWithContext(ctx context.Context, maxSize int) (*Semaphore, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Semaphore{
		c:      make(chan struct{}, maxSize),
		cancel: cancel,
	}, ctx
}

func (s *Semaphore) Go(f func() error) {
	s.wg.Add(1)
	s.c <- struct{}{}
	go func() {
		defer func() {
			if err := recover(); err != nil {

			}
		}()
		defer func() {
			<-s.c
			s.wg.Done()
		}()
		if err := f(); err != nil {
			s.errOnce.Do(func() {
				s.err = err
				if s.cancel != nil {
					s.cancel()
				}
			})
		}
	}()
}

func (s *Semaphore) Wait() error {
	s.wg.Wait()
	if s.cancel != nil {
		s.cancel()
	}
	return s.err
}
