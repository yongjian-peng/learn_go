package main

import "fmt"

func main() {
	var s []int
	for i := 0; i < 3; i++ {
		s = append(s, i)
	}

	modifySlice2(s)
	fmt.Println(s)
}

func modifySlice2(s []int) {
	s = append(s, 2048) // 传的值的类型是不一样的，发生了变化
	s[0] = 1024
}
