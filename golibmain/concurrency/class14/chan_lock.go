package main

import (
	"fmt"
	"time"
)

// 使用chan实现互斥锁
type Mutex struct {
	ch chan struct{}
}

// 使用锁需要初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

// 请求锁，直到获取
func (m Mutex) Lock() {
	<-m.ch
}

// 解锁
func (m *Mutex) UnLock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")

	}
}

// 尝试获得锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:

	}
	return false
}

// 加入一个超时的设置
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}

// 使用chan实现互斥锁，至少有两种方式。
// 一种方式是先初始化一个capacity等于1的Channel,
// 然后再放入一个元素，这个元素就代表锁，谁获取了这个元素，就等于获取了这把锁。
// 另一种方式是，先初始化一个capacity等于1的Channel,它的“空槽”代表锁，谁能成功地把元素发送到这个Channel，谁就获取了这把锁。
func main() {
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
}
