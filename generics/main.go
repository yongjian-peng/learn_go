package main

import (
	// "fmt"
	// "strconv"
	"./hao"
)

type Number interface {
	int64 | float64
}

func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
    var s V
    for _, v := range m {
        s += v
    }
    return s
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	// $ints = 
	// ints := map[string]int64{
	// 	"first":  34,
	// 	"second": 12,
	// }

	// floats := map[string]float64{
	// 	"first":  35.98,
	// 	"second": 26.99,
	// }
	// fmt.Printf("Non-Generic Sums: %v and %v\n",
	// 	SumInts(ints),
	// 	SumFloats(floats))

	// fmt.Printf("Generic Sums: %v and %v\n",
	// 	SumIntsOrFloats[string, int64](ints),
	// 	SumIntsOrFloats[string, float64](floats))

	// fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
	// 	SumIntsOrFloats(ints),
	// 	SumIntsOrFloats(floats))

	// fmt.Printf("Generic Sums with Constraint: %v and %v\n",
	// 	SumNumbers(ints),
	// 	SumNumbers(floats))
	// 数组切片
	// hao.RunSlice();
	// 字符串切片
	// hao.RunSlice2()
	// map 深度比较是否相等
	// hao.RunSlice3()
	// 结构体 打印
	// hao.RunSlice4()
	// 接口编程模式
	// hao.RunSlice5()
	// 接口验证
	// hao.RunSlice6();

	// go 资源清理
	// hao.RunSlice7();

	// 错误处理方式
	// hao.RunSlice8()

	// 指针 point 学习
	hao.RunPoint()

	// 函数式编程
	// hao.RunFunctional()
}
