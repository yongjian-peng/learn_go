package main

import "fmt"

type person struct {
	age  int
	name string
}

func modify(p person) {
	p.age = p.age + 10
}

// 引用的例程
func main() {
	jim := person{
		age:  10,
		name: "jim",
	}
	fmt.Println("jim", jim)

	modify(jim)
	fmt.Println("jim2", jim) // 是不会改变 person 的值 值拷贝
}
