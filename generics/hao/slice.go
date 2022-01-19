package hao

import (
	"fmt"
	"unsafe"
)

type slice struct {
	array unsafe.Pointer // 指向存放数据的数组指针
	len   int            // 长度有多大
	cap   int            // 容量有多大
}

func RunSlice() {
	foo = make([]int, 5)
	foo[3] = 42
	foo[4] = 100

	bar := foo[1:4]
	bar[1] = 99
	fmt.Println(bar)
}
