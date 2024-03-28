package main

import "fmt"

// 给定一个排序数组和目标值， 在数组中找到目标值， 并返回其索引，如果目标值不存在，则返回数组被按顺序插入的位置。
func main() {
	nums := []int{1, 3, 5, 8, 9}

	length := len(nums)

	index := search(nums, 0, length-1, 6)

	fmt.Println("index:", index)

	fmt.Println("nums=>", nums)

	left := nums[:index]

	fmt.Println("left:", left)
	right := nums[index:]

	newNums := make([]int, length+1)
	newNums = append(newNums, left...)
	newNums = append(newNums, 6)
	newNums = append(newNums, right...)

	fmt.Println("nums:", nums)
	fmt.Println("newNums:", newNums)

	fmt.Println("right:", right)

}

// search 二分查找，如果不存在，返回应该插入的值。
func search(nums []int, left, right, target int) int {
	for left <= right {
		var mid = left + (right-left)/2
		fmt.Println("left - in ", left)
		if target == nums[mid] {
			return mid
		} else if nums[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	fmt.Println("left", left)

	return left
}
