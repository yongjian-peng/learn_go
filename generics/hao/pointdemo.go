package hao

import "fmt"

func RunPoint() {
	// b := 255
	// var a *int = &b // a 是int 指针 指向b的内存地址
	// fmt.Printf("Type of is: %T\n", a)
	// fmt.Printf("address of b is", a)

	// a := 22
	// var b *int // b 这个变量是int类型的指针变量，变量的值，也只能是 int 型指针
	// if b == nil {
	// 	// 指针的零值是nil
	// 	fmt.Println("b is", b)
	// 	b = &a
	// 	fmt.Println("b after initialization is", b)
	// }

	// b := 255
	// a := &b
	// fmt.Println("address of b is", a) // 打印b的内存地址
	// fmt.Println("value of b is", *a)  // 打印b的值 可以通过*a指针
	// fmt.Println("value of b is", &a)  // 打印a的内存地址 可以通过a指针

	// b := 255
	// a := &b
	// fmt.Println("address of b is:", a) // b的内存地址
	// fmt.Println("value of b is:", *a)
	// *a++ // 通过a的指针加一
	// fmt.Println("new value of b is:", b)

	// a := 58
	// fmt.Println("value of a before func call is:", a)
	// b := &a
	// change(b) // 指针变量b 改变a的值 a=55
	// fmt.Println("value of a after call is:", a)

	/**
	不要将指向数组的指针，作为参数传递给函数，改用切片
	*/
	// a := [3]int{89, 90, 91}
	// modify(&a)
	// fmt.Println(a)

	a := [3]int{89, 90, 91}
	modify(a[:]) // 传递a的一个切片
	fmt.Println(a)

}

// 函数传递指针，改变参数的内存地址
func change(val *int) {
	*val = 55
}

// 修改数组的值 传递指向数组的指针 作为参数 并且对其修改
// func modify(arr *[3]int) {
// 	// (*arr)[0] = 90
// 	arr[0] = 90 // 也可以这么写 这是上面的简写形式(*arr)[x] 可以写成a[x]
// }

// 切片方式修改函数 这个方法是修改函数最常用的 最好的方法
func modify(sls []int) {
	sls[0] = 91
}
