package main

var a, b int

func f() {
	a = 1 // w之前的写操作
	b = 2 // 写操作w
}

func g() {
	print(b) // w之前的写操作
	print(a) // ???
}

func main() {
	go f() // g1
	g()    // g2
}
