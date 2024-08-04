package main

import "fmt"

// reduce
// 写一个程序， reduce 函数就是把map处理的结果相乘

func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil { // 异常检查
		return nil
	}
	out := <-in         // 先读取第一个元素
	for v := range in { // 实现reduce的主要逻辑
		out = fn(out, v)
	}
	return out
}

// 生成一个数据流
func asStreamMR(done <-chan struct{}) <-chan interface{} {
	s := make(chan interface{})
	values := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(s)
		for _, v := range values { // 从数组生成
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()
	return s
}

func main() {
	p := make(chan struct{}, 1)
	in := asStreamMR(p)

	// reduce操作，对map的结果进行累加
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) * v.(int)
	}

	sum := reduce(in, reduceFn) // 返回累加结果
	fmt.Println("sum=>", sum)
	p <- struct{}{}

}
