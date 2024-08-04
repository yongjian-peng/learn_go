package main

import (
	"fmt"
	"unsafe"
)

type Num struct {
	i string
	j int64
}

type Sample struct {
	BoolValue  bool
	FloatValue float64
}

// 引用 https://segmentfault.com/a/1190000017389782
func main() {
	n := Num{
		i: "EDDYCJY",
		j: 1,
	}

	nPointer := unsafe.Pointer(&n)
	fmt.Println("nPointer=>", nPointer)

	niPointer := (*string)(nPointer)
	*niPointer = "煎鱼"

	njPointer := (*int64)(unsafe.Pointer(uintptr(nPointer) + unsafe.Offsetof(n.j)))

	*njPointer = 2
	fmt.Printf("n.i: %s, n.j: %d\n", n.i, n.j)

	s := Sample{BoolValue: true, FloatValue: 3.14}

	fmt.Println("BoolValue Align:", unsafe.Alignof(s.BoolValue))
	fmt.Println("FloatValue Align:", unsafe.Alignof(s.FloatValue))
	fmt.Println("BoolValue Offset:", unsafe.Offsetof(s.BoolValue))
	fmt.Println("FloatValue Offset:", unsafe.Offsetof(s.FloatValue))
	fmt.Println("BoolValue Size:", unsafe.Sizeof(s.BoolValue))
	fmt.Println("FloatValue Size:", unsafe.Sizeof(s.FloatValue))
}
