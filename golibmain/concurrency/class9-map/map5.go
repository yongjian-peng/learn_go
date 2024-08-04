package main

import (
	"fmt"
	"sync"
)

func main() {
	// 创建一个新的 sync.Map
	syncMap := &sync.Map{}

	// 设置键值对
	syncMap.Store("key1", "value1")
	syncMap.Store("key2", "value2")
	syncMap.Store("key3", "value3")

	// 从map中读取值
	value, exists := syncMap.Load("key2")
	if exists {
		fmt.Println("value=>", value)
	}
	if !exists {
		fmt.Println("Key not found")
	}

	// 删除键
	syncMap.Delete("key3")

	// 遍历map
	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // 继续遍历
	})
}
