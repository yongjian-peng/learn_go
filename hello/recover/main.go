package main

import (
	"fmt"
	"time"
)

func write() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("catch err")
		}
	}()

	var m map[string]int
	m["abc"] = 1
}

func main() {
	go write()
	time.Sleep(time.Second * 2)
}
