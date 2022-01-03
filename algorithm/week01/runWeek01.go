package week01

import (
	"fmt"
)

func RunWeek01() {
	fmt.Println("RunWeek01")
	// 柱状图中最大矩形
	largest_rect := []int{2, 1, 5, 6, 2, 3}
	ans := LargestRectangleArea(largest_rect)

	fmt.Println(ans)
}
