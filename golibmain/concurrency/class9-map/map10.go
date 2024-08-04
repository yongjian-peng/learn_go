package main

import "fmt"

func main() {
	m := make(map[int]int)
	m[1] = 1

	fmt.Println(&m[1]) // map 存的是值，会发生copy, map 的 key value 都不可以取地址，不建议放大的map
}
