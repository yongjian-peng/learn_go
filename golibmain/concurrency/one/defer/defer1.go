package main

import "fmt"

var a, b, c int

const s1 = "chain"

func main() {
	a = 1
	b = 2
	c = 3

	fmt.Println(a, b, c, s1)

	s := make([]int, 5)

	s = append(s, 1, 2, 3)

	fmt.Println(s) // [0 0 0 0 0 1 2 3]

	//deferCall()
}

func deferCall() {
	/**
	打印后
	打印中
	打印前
	panic: 触发异常
	*/
	defer func() {
		fmt.Println("打印前")
		/*if err := recover(); err != nil {
			fmt.Println(err)
		}*/
	}()

	defer func() { fmt.Println("打印中") }()

	defer func() { fmt.Println("打印后") }()

	panic("触发异常") // 最后才会显示这个打印
}
