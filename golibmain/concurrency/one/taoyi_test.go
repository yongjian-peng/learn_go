package main

import "testing"

// go test -bench=BenchmarkSlice .\taoyi_test.go
// sliceEscape 发生逃逸，在堆上申请切片
func sliceEscape() {
	number := 10
	s1 := make([]int, 0, number)
	for i := 0; i < number; i++ {
		s1 = append(s1, i)
	}
}

// sliceNoEscape 不逃逸，限制在栈上
func sliceNoEscape() {
	s1 := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		s1 = append(s1, i)
	}
}

func BenchmarkSliceEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sliceEscape()
	}
}

func BenchmarkSliceNoEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sliceNoEscape()
	}
}
