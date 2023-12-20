package main

import "fmt"

// map[key]函数返回结果可以是一个值，也可以是两个值
func main() {
	var m = make(map[string]int)
	m["a"] = 0
	fmt.Printf("a=%d; b=%d\n", m["a"], m["b"])

	av, aexisted := m["a"]
	bv, bexisted := m["b"]
	fmt.Printf("a=%d, existed: %t; b=%d, existed: %t\n", av, aexisted, bv, bexisted)
}
