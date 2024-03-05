package main

// 内存逃逸 例程 编译命令 go build -gcflags="-m -m -l" taoyi.go
func main() {
	number := 10
	s1 := make([]int, 0, number)
	for i := 0; i < number; i++ {
		s1 = append(s1, i)
	}

	s2 := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		s2 = append(s2, i)
	}
}

/**
.\taoyi.go:6:12: make([]int, 0, number) escapes to heap:
.\taoyi.go:6:12:   flow: {heap} = &{storage for make([]int, 0, number)}:
.\taoyi.go:6:12:     from make([]int, 0, number) (non-constant size) at .\taoyi.go:6:12
.\taoyi.go:6:12: make([]int, 0, number) escapes to heap
.\taoyi.go:11:12: make([]int, 0, 10) does not escape
*/
