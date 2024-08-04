package main

import (
	"fmt"
	"sync"
	"time"
)

// MuOnce 一个组合的并发源语
type MuOnce struct {
	sync.RWMutex
	sync.Once
	mtime time.Time
	vals  []string
}

// refresh 相当于reset方法，会将m.Once重新复制一个Once
func (m *MuOnce) refresh() {
	m.Lock()
	defer m.Unlock()
	m.Once = sync.Once{}
	m.mtime = time.Now()
	m.vals = []string{m.mtime.String()}
}

// strings 获取某个初始化的值，如果超过某个时间，会reset Once
func (m *MuOnce) strings() []string {
	now := time.Now()
	m.RLock()
	if now.After((m.mtime)) {
		defer m.Do(m.refresh) // 使用refresh函数重新初始化
	}
	vals := m.vals
	m.RUnlock()
	return vals
}

func main() {
	fmt.Println("Hello playground")

	m := new(MuOnce)
	fmt.Println("m.strings()=>", m.strings())
	fmt.Println("m.strings()=>", m.strings())
}
