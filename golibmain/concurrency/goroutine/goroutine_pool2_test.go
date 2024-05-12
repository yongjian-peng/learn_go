package main

import (
	"github.com/bmizerany/assert"
	"math"
	"testing"
	"time"
)

func sleeps() error {
	time.Sleep(time.Second)
	return nil
}

func TestSemaphore(t *testing.T) {
	// 最大并发 => 执行任务数量
	sema := NewSemaphore(4)
	now := time.Now()
	for i := 0; i < 4; i++ {
		sema.Go(sleeps)
	}

	err := sema.Wait()
	assert.Equal(t, nil, err)
	sec := math.Round(time.Since(now).Seconds())
	assert.Equal(t, 1, int(sec))
}
