package main

import "fmt"

type mapKey struct {
	key int
}

// key 类型的 K 必须是可比较的
// 在 Go语言中， bool,整数，浮点数，复数，字符串，指针，Channel，都是可以比较的，包含可以比较的元素的 struct和数组
// 而 slice、map、函数值都是吧可以比较的
func main() {
	var m = make(map[mapKey]string)
	var key = mapKey{10}

	m[key] = "hello"
	fmt.Printf("m[key]=%s\n", m[key])

	// 修改key的字段的值后再次查询map 无法获取刚才add进去的值
	key.key = 100
	fmt.Printf("再次查询m[key]=%s\n", m[key])
}
