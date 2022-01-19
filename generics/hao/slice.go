package hao

import (
	"bytes"
	"fmt"
	"unsafe"
)

type slice struct {
	array unsafe.Pointer // 指向存放数据的数组指针
	len   int            // 长度有多大
	cap   int            // 容量有多大
}

func RunSlice() {
	fmt.Println("bar")
	a := make([]int, 32)
	a[5] = 5
	b := a[1:16]
	a = append(a, 1)
	a[2] = 42
	fmt.Println(a)
	fmt.Println(b)
}

func RunSlice2() {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]

	fmt.Println("dir1 =>", string(dir1))
	fmt.Println("dir2 =>", string(dir2))
}
