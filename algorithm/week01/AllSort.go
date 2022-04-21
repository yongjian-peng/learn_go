package week01

import "fmt"

func RunSelectSort() {

	fmt.Println("OK")

	arr := []int{4, 6, 9, 10, 33}

	// SelectSort(arr)

	InsertSort(arr)

	fmt.Println(arr)

}

func SelectSort(arr []int) {
	// 选择排序 再循环中 找到最小的那个 压入到新的数组中 交换

	lenght := len(arr)
	if lenght <= 1 {
		return
	}

	for i := 0; i < lenght; i++ {
		curr := i

		for j := i + 1; j < lenght; j++ {
			if arr[j] <= arr[curr] {
				curr = j
			}
		}

		if curr != i {
			arr[i], arr[curr] = arr[curr], arr[i]
		}
	}
}

func InsertSort(arr []int) {
	lenght := len(arr)
	if lenght <= 1 {
		return
	}
	// 插入排序 和当前的元素 做对比 循环交换位置
	for i := 1; i < lenght; i++ {
		curr := arr[i]
		for j := i - 1; j >= 0; j-- {
			if arr[j] >= curr {
				arr[j+1] = arr[j]
				arr[j] = curr
			} else {
				break
			}
		}

	}
}
