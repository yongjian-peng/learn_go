package main

import (
	"fmt"
	"sync"
)

func printNumbers(ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 28; i += 2 {
		<-ch // Wait for signal from the other goroutine
		fmt.Print(i, i+1)
		ch <- true // Signal the other goroutine to proceed
	}
}

func printLetters(ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	index := 0

	for i := 0; i < 28; i += 2 {
		<-ch // Wait for signal from the other goroutine
		fmt.Print(letters[index : index+2])
		index = (index + 2) % len(letters) // Wrap around to the beginning if we reach the end
		ch <- true                         // Signal the other goroutine to proceed
	}
}

/**
1.交替使用数字和字母。
使用两个程序交替打印序列，一个程序打印数字，另一个程序
打印结果如下:12ab34cd56ef78gh910j1112kl1314mn1516op1718qr1920st2122uv2324wx2526yz2728gooutine打印字母

*/
func main() {
	ch := make(chan bool, 1) // Buffered channel with capacity 1
	var wg sync.WaitGroup
	wg.Add(2)

	go printLetters(ch, &wg)
	go printNumbers(ch, &wg)

	ch <- true // Start the first goroutine

	wg.Wait() // Wait for both goroutines to finish
	fmt.Println()
}
