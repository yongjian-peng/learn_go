package week01

import "fmt"

func RunTrap() {
	heights := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}

	ant := trap(heights)

	fmt.Println(ant)
}

func trap(height []int) int {
	var ant int
	var stack []Rect
	for i := 0; i < len(height); i++ {
		var width = 0
		for len(stack) > 0 && stack[len(stack)-1].height <= height[i] {
			width += stack[len(stack)-1].width
			var bottom = stack[len(stack)-1].height
			stack = stack[0:(len(stack) - 1)]
			if len(stack) == 0 {
				continue
			}
			var up = min(height[i], stack[len(stack)-1].height)
			ant += width * (up - bottom)
		}

		stack = append(stack, Rect{width: width + 1, height: height[i]})
	}

	return ant
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
