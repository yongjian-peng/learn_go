package hao

import (
	"bytes"
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
<<<<<<< HEAD
	fmt.Println("bar")
	a := make([]int, 32)
	a[5] = 5
	b := a[1:16]
	a = append(a, 1)
	a[2] = 42
	fmt.Println(a)
	fmt.Println(b)
}

func RunSlice2() {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]
=======
	var foo []int
	foo = make([]int, 5)
	foo[3] = 42
	foo[4] = 100
>>>>>>> aef7066e6f57e8307f674d1ed48448c7511f36f7

	fmt.Println("dir1 =>", string(dir1))
	fmt.Println("dir2 =>", string(dir2))
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
