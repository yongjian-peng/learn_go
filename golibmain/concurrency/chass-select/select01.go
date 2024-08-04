package main

import "fmt"

func main() {
	s := make([]int, 0, 10000)
	for i := 0; i < 10000; i++ {
		s = append(s, 1)
	}

	s2 := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		s2[i] = i
	}
	caseone()
}

func caseone() {
	var s []int
	for i := 0; i < 3; i++ {
		s = append(s, i)
	}

	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(s []int) {
	s[0] = 1024
}
