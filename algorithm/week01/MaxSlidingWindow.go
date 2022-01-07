package week01

import "fmt"

func RunMaxSlidingWindow() {
	nums := []int{1, 3, -1, -3, 5, 3, 6, 7}

	ant := MaxSlidingWindow(nums, 3)

	fmt.Println(ant)

}

func MaxSlidingWindow(nums []int, k int) []int {
	if k == 1 {
		return nums
	}
	var ant []int
	var queue []int
	for i := 0; i < len(nums); i++ {
		for len(queue) > 0 && queue[0] <= i-k {
			queue = queue[1:]
		}
		for len(queue) > 0 && nums[queue[len(queue)-1]] <= nums[i] {
			queue = queue[0:(len(queue) - 1)]
		}
		queue = append(queue, i)
		if i >= k-1 {
			ant = append(ant, nums[queue[0]])
		}
	}
	return ant
}
