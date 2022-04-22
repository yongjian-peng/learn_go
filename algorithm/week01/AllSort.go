package week01

import "fmt"

func RunSelectSort() {

	fmt.Println("OK")

	arr := []int{4, 6, 19, 10, 33, 7}

	// SelectSort(arr)

	// InsertSort(arr)

	// 归并排序
	lenght := len(arr)
	MergeSort(arr, 0, lenght-1)

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

func MergeSort(arr []int, l int, r int) {
	// 有一个二分
	if l >= r {
		return
	}

	mid := (l + r) >> 1

	MergeSort(arr, 0, mid)
	MergeSort(arr, mid+1, r)

	MergeArray(arr, l, mid, r)

	// 合并两个数组 返回一个数组
}

func MergeArray(arr []int, left int, mid int, right int) {
	len := right - left + 1
	var temp []int

	i := left
	j := mid + 1

	for k := 0; k < len; k++ {
		if j > right || (i <= mid && arr[i] <= arr[j]) {
			temp = append(temp, arr[i])
			i++
		} else {
			temp = append(temp, arr[j])
			j++
		}
	}

	for k := 0; k < len; k++ {
		arr[left+k] = temp[k]
	}
}
