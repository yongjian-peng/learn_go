package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan string, 1)

	go func() {
		time.Sleep(3 * time.Second)
		c <- "aa"
	}()

	s, ok := <-c
	if ok {
		fmt.Println("ok")
	}

	fmt.Println("s=>", s)

	fmt.Println("chan")
}
