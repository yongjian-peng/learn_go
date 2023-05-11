package main

import "sync"

type T int

type MyChannel struct {
	C      chan T
	closed bool
	mutex  sync.Mutex
}

// 使用 sync.Mutex 来避免多次关闭一个通道
// 出自：https://go101.org/article/channel-closing.html
func NewMyChannel() *MyChannel {
	return &MyChannel{C: make(chan T)}
}

func (mc *MyChannel) SafeClose() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if !mc.closed {
		close(mc.C)
		mc.closed = true
	}
}

func (mc *MyChannel) IsClosed() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.closed
}
