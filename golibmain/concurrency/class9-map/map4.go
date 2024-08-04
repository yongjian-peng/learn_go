package main

import (
	"fmt"
	"sync"
)

// RWMap 一个读写锁保护的线程安全的map
type RWMap struct {
	sync.RWMutex             // 读写锁保护下面的map字段
	m            map[int]int // 保护的map
}

// NewRWMap 新建一个RWMap
func NewRWMap(n int) *RWMap {
	return &RWMap{
		m: make(map[int]int, n),
	}
}

// Get 从map中读取一个值
func (m *RWMap) Get(k int) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	v, existed := m.m[k] // 在锁的保护下从map中读取
	return v, existed
}

// Set 设置一个键值对
func (m *RWMap) Set(k int, v int) {
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

// Delete 删除一个键
func (m *RWMap) Delete(k int) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, k)
}

// Len 获取map的长度
func (m *RWMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

// Each 遍历map
func (m *RWMap) Each(f func(k, v int) bool) {
	m.RLock() // 遍历期间一致持有读锁
	defer m.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
}

func main() {
	// TODO 完成调用demo，加深理解
	newMap := NewRWMap(10)
	newMap.Set(1, 20)
	newMap.Set(2, 200)
	newMap.Set(3, 300)

	one, _ := newMap.Get(1)

	fmt.Println("one=>", one)

	fmt.Println("len", newMap.Len())
	// TODO each调用
	newMap.Each(func(k, v int) bool {
		fmt.Printf("Key: %d, Value: %d\n", k, v)
		return true // 继续遍历
	})
}

// ConcurrentMap 并发map实例，
// 引用：https://github.com/orcaman/concurrent-map/blob/master/README.md
