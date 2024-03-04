package main

import (
	"fmt"
	"sync"
	"time"
)

// FanOut fan-out pattern implementation
func FanOut(input <-chan int, outputs []chan<- int, async bool) {
	var wg sync.WaitGroup

	for _, out := range outputs {
		wg.Add(1)
		go func(out chan<- int) {
			defer wg.Done()
			for value := range input {
				if async {
					// Asynchronous send
					go func(value int) {
						select {
						case out <- value:
						case <-time.After(time.Second):
							// Output channel might be closed, just drop the value
						}
					}(value)
				} else {
					// Synchronous send
					select {
					case out <- value:
					case <-time.After(time.Second):
						// Output channel might be closed, just drop the value
					}
				}
			}
		}(out)
	}

	// Wait for all output goroutines to finish
	wg.Wait()
}

// 扇出模式 一个输入，多个输出
func main() {
	// Create an input channel
	input := make(chan int)

	// Create output channels
	outputs := []chan<- int{
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
	}

	// Start fan-out goroutine with asynchronous mode
	go FanOut(input, outputs, false)

	// Send values to the input channel
	for i := 0; i < 10; i++ {
		input <- i
		time.Sleep(time.Millisecond * 100) // Simulate work
	}

	// Close the input channel to signal the end
	close(input)

	// Wait for the fan-out goroutine to finish
	time.Sleep(time.Second) // Wait for asynchronous sends to complete

	// Close the output channels
	for _, out := range outputs {
		close(out)
	}

	// Consume the output channels
	for i, out := range outputs {
		fmt.Printf("Output channel %d:\n", i+1)
		//for value := range out {
		//	fmt.Printf("%d ", out)
		//}
		fmt.Println(out)
	}
}
