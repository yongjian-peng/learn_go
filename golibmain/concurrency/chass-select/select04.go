package main

import "fmt"

// 如果没有发生扩容，修噶在原来的内存中
// 如果发生了扩容，修改会在新的内存中
func main() {
	var s []int
	for i := 0; i < 3; i++ {
		s = append(s, i)
	}

	modifySlice4(s)
	fmt.Println(s) // 1024 1 2
}

func modifySlice4(s []int) {
	s[0] = 1024
	s = append(s, 2048)
	s = append(s, 4096)

}
