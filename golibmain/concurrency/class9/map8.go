package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	NullRes()
	ArrayRes()

}

func NullRes() {
	var s []int

	b, _ := json.Marshal(s)
	fmt.Printf(string(b)) // return null
}

func ArrayRes() {
	// 这是编辑器的语法检查的问题
	var s []int{}

	b, _ := json.Marshal(s)
	fmt.Printf(string(b)) // return []
}
