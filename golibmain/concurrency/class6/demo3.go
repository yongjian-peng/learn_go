package main

import "fmt"

func main() {
	c := make(chan string)

	go func() {
		c <- "aa"
	}()

	//s, ok := <-c
	//if ok {
	//	fmt.Println("ok")
	//}

	fmt.Println("chan")
}
