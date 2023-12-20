package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once

	// 第一个初始化函数
	f1 := func() {
		fmt.Println("第一次打印")
	}

	once.Do(f1)

	// 第二个初始化函数
	f2 := func() {
		fmt.Println("第二次打印")
	}

	once.Do(f2)
}
