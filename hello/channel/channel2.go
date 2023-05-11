package main

import "fmt"

type T int

// 出自：https://go101.org/article/channel-closing.html
func IsClosed(ch <-chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

// 检查 channel 这不是是否关闭的通用方法
func main() {
	c := make(chan T)
	fmt.Println((IsClosed(c)))
	close(c)
	fmt.Println(IsClosed(c))
}
