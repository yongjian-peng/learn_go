package main

import "fmt"

// map-reduce
// 写一个程序，使用map-reduce模式处理一组整数，map函数就是为每个整数乘以10， reduce 函数就是把map处理的结果累加起来

func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{}) // 创建一个输出chan
	if in == nil {                // 异常检查
		close(out)
		return out
	}

	go func() { // 启动一个goroutine,实现map的主要逻辑
		defer close(out)
		for v := range in { // 从输入chan读取数据，执行业务操作，也就是map操作
			out <- fn(v)
		}
	}()
	return out
}

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

	// map操作，乘以10
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}

	// reduce操作，对map的结果进行累加
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) + v.(int)
	}

	sum := reduce(mapChan(in, mapFn), reduceFn) // 返回累加结果
	fmt.Println("sum=>", sum)
	p <- struct{}{}

}
