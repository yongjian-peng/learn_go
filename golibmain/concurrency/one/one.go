package main

import "fmt"

const (
	mutexLocke = 1 << iota
	mutexWoke
	mutexStart
	mutexWaiterShi = iota
)

func main() {

	var waitStartTime int64

	fmt.Println("waitStartTime=>", waitStartTime)

	fmt.Println("mutexLocke=>", mutexLocke)
	fmt.Println("mutexWoke=>", mutexWoke)
	fmt.Println("mutexStart=>", mutexStart)
	fmt.Println("mutexWaiterShi=>", mutexWaiterShi)

	old := 1
	fmt.Printf("old,%b\n", old)

	new := old + 1<<mutexWaiterShi
	fmt.Println("new=>", new)
	fmt.Printf("new,%b\n", new)

	new = new + 1<<mutexWaiterShi
	fmt.Println("new2=>", new)
	fmt.Printf("new2,%b\n", new)

	new = new + 1<<mutexWaiterShi
	fmt.Println("new3=>", new)
	fmt.Printf("new3,%b\n", new)

}
