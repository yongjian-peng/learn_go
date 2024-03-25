package main

import "fmt"

func merge(left, right []int) []int {
	result := make([]int, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result[i+j] = left[i]
			i++
		} else {
			result[i+j] = right[j]
			j++
		}
	}

	for i < len(left) {
		result[i+j] = left[i]
		i++
	}

	for j < len(right) {
		result[i+j] = right[j]
		j++
	}

	return result
}

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

func main() {
	arr := []int{38, 27, 43, 3, 9, 82, 10}
	sortedArr := mergeSort(arr)
	fmt.Println("Sorted array is:", sortedArr)
}
