package main

import "fmt"

const (
	a = 1 << iota // a == 1  (iota == 0)
	b = 1 << iota // b == 2  (iota == 1)
	c = 3         // c == 3  (iota == 2, unused)
	d = 1 << iota // d == 8  (iota == 3)
)

// var s interface{ comparable }

func main() {

	// 	[...]Point{{1.5, -3.5}, {0, 0}}     // same as [...]Point{Point{1.5, -3.5}, Point{0, 0}}
	// [][]int{{1, 2, 3}, {4, 5}}          // same as [][]int{[]int{1, 2, 3}, []int{4, 5}}
	// [][]Point{{{0, 1}, {1, 2}}}         // same as [][]Point{[]Point{Point{0, 1}, Point{1, 2}}}
	// map[string]Point{"orig": {0, 0}}    // same as map[string]Point{"orig": Point{0, 0}}
	// map[Point]string{{0, 0}: "orig"}    // same as map[Point]string{Point{0, 0}: "orig"}

	// type PPoint *Point
	// [2]*Point{{1.5, -3.5}, {}}          // same as [2]*Point{&Point{1.5, -3.5}, &Point{}}
	// [2]PPoint{{1.5, -3.5}, {}}          // same as [2]PPoint{PPoint(&Point{1.5, -3.5}), PPoint(&Point{})}

	// list of prime numbers
	// primes := []int{2, 3, 5, 7, 9, 2147483647}

	// 主要表达
	// // vowels[ch] is true if ch is a vowel
	// vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}

	// // the array [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}
	// filter := [10]float32{-1, 4: -0.1, -0.1, 9: -1}

	// // frequencies in Hz for equal-tempered scale (A4 = 440Hz)
	// noteFrequency := map[string]float32{
	// 	"C0": 16.35, "D0": 18.35, "E0": 20.60, "F0": 21.83,
	// 	"G0": 24.50, "A0": 27.50, "B0": 30.87,

	var re, im = complexSqrt(-1)
	fmt.Println(re, im)
}
