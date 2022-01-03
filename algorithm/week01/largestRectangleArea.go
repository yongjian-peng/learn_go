package week01

import "fmt"

// 柱张图中的最大矩形 单调栈解法
type Rect struct {
	width  int
	height int
}

// import "fmt"

func RunLargestRectangleArea() {
	// 柱状图中最大矩形
	largest_rect := []int{2, 1, 5, 6, 2, 3}
	ans := LargestRectangleArea(largest_rect)

	fmt.Println(ans)
}

func LargestRectangleArea(heights []int) int {
	//fmt.Println("largestRectangleArea")
	var ans int
	var stack []*Rect
	heights = append(heights, 0)
	for i := 0; i < len(heights); i++ {
		var width int
		for len(stack) > 0 && stack[len(stack)-1].height >= heights[i] {
			width += stack[len(stack)-1].width
			if ans < stack[len(stack)-1].height*width {
				ans = stack[len(stack)-1].height * width
			}
			stack = stack[0:(len(stack) - 1)]
		}
		stack = append(stack, &Rect{width: width + 1, height: heights[i]})
	}
	return ans
}
