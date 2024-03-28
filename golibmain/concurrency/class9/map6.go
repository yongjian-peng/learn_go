package main

import (
	"fmt"
	"time"
)

func main() {
	s := make(map[int]int)
	for i := 0; i < 100; i++ {
		go func(i int) {
			s[i] = i
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Printf("map%d个元素是%d", i, s[i])
		}(i)
	}
	time.Sleep(1 * time.Second)
}
