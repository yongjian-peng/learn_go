package week01

import "fmt"

func RunMaxSlidingWindow() {
	nums := []int{1, 3, -1, -3, 5, 3, 6, 7}

	ant := MaxSlidingWindow2(nums, 3)

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

// queue 中不存在位置，直接存值，减少nums寻址次数
func MaxSlidingWindow2(nums []int, k int) []int {
	if k == 1 {
		return nums
	}
	var ant []int
	var queue []int
	for i := 0; i < len(nums); i++ {
		for len(queue) > 0 && queue[0] <= i-k {
			queue = queue[1:]
		}
		for len(queue) > 0 && nums[queue[len(queue)-1]] < nums[i] {
			queue = queue[0:(len(queue) - 1)]
		}
		queue = append(queue, i)
		if i >= k-1 {
			ant = append(ant, nums[queue[0]])
		}
	}
	return ant
}
