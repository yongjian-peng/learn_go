package main

import "fmt"

// 接受返回值的时候，第二个值是bool类型
// 代表是否成功的从 chan 中读取到一个值，如果第二个参数是false,
// chan 已经被 close 而且 chan 中没有缓存的数据。
// 这个时候，第一个值是零值，所以，如果从 chan 读取到一个零值，可能是 sender 真正的发送零值。
// 也可能是 closed 的并且没有缓存元素产生的零值
func main() {
	// 创建一个带有缓冲区的通道
	ch := make(chan int, 1)

	// 发送数据到通道
	ch <- 0

	// 从通道接受的值
	value, ok := <-ch

	// 检查是否成功的接受到值
	if ok {
		fmt.Println("接受到的值：", value)
	} else {
		fmt.Println("未接受到值")
	}
}
