package main

import (
	"fmt"
	"sync"
)

func fooReen(l sync.Locker) {
	fmt.Println("in fooReen")
	l.Lock()
	bar(l)
	l.Unlock()
}

func bar(l sync.Locker) {
	l.Lock()
	fmt.Println("in bar")
	l.Unlock()
}

func main() {
	l := &sync.Mutex{}
	fooReen(l)
}
