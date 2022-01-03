package week01

// 柱张图中的最大矩形 单调栈解法
type Rect struct {
	width  int
	height int
}

// import "fmt"

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
