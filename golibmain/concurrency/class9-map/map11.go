package main

import (
	"fmt"
	"runtime"
)

func main() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Println(mem.HeapInuse)
	m := make(map[int]int)
	for i := 0; i < 10000; i++ {
		m[i] = i
	}

	runtime.ReadMemStats(&mem)
	fmt.Println(mem.HeapInuse)
	for k := range m {
		delete(m, k)
	}

	runtime.ReadMemStats(&mem)
	fmt.Println(mem.HeapInuse)
}
