package main

import "fmt"

func main() {
	for i := 0; i < 8; i++ {
		fmt.Println("i=>", i%4)
	}
}
