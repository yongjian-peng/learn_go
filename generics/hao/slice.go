package hao

import (
	"fmt"
	"reflect"
	"unsafe"
)

type slice struct {
	array unsafe.Pointer // 指向存放数据的数组指针
	len   int            // 长度有多大
	cap   int            // 容量有多大
}

type Person struct {
	Name   string // 指向存放数据的数组指针
	Sexual string // 长度有多大
	Age    int    // 容量有多大
}

func RunSlice() {
	var foo []int
	foo = make([]int, 5)
	foo[3] = 42
	foo[4] = 100

	bar := foo[1:4]
	bar[1] = 99
	fmt.Println(bar)
}

func RunSlice3() {
	// var data []int{}
	// v1 := data{}
	// v2 := data{}

	// fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))
}

func PrintPerson(p *Person) {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
		p.Name, p.Sexual, p.Age)
}
func (p *Person) Print() {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
		p.Name, p.Sexual, p.Age)
}

func RunSlice4() {
	var p = Person{
		Name:   "Hao Chen",
		Sexual: "Male",
		Age:    44,
	}
	PrintPerson(&p)
	p.Print()
}
