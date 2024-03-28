package main

import "fmt"

func main() {

	//for i := 0; i < 100; i++ {
	//	err01()
	//}
	//err02()
	err03()
}

func err01() {
	i := 1

	defer func() {
		i++
	}()

	fmt.Println("i:", i)
}

func err02() {
	i := 1

	defer add(i)

	fmt.Println("i:", i)
}

func add(i int) {
	i++
}

func err03() {
	done := make(chan bool)

	values := []string{"a", "b", "c"}

	for _, v := range values {
		go func() {
			fmt.Println(v)
			//done <- true
			<-done
		}()
	}

	for _ = range values {
		done <- true
		//<-done
	}
}
