package main

import "fmt"

func main() {
	var s []int
	for i := 0; i < 3; i++ {
		s = append(s, i)
	}

	modifySlice3(s)
	fmt.Println(s) // 0 1 2
}

func modifySlice3(s []int) {
	s = append(s, 2048)
	s = append(s, 4096) // 扩容了，重新分配了内存。
	s[0] = 1024
}
