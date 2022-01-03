package week01

import "fmt"

func RunMergeTwoArray() {
	nums1 := []int{1, 2, 3, 0, 0, 0}

	nums2 := []int{2, 5, 6}

	var m, n = 3, 3

	// fmt.Println(m)

	MergeTwoArray(nums1, m, nums2, n)
	fmt.Println(nums1)
}

func MergeTwoArray(nums1 []int, m int, nums2 []int, n int) {
	// 合并两个数组 倒序排序 可以节省一个数组的拷贝
	i := m - 1
	j := n - 1
	for k := m + n - 1; k >= 0; k-- {
		// 判断是否越界
		if j < 0 || (i >= 0 && nums1[i] >= nums2[j]) {
			nums1[k] = nums1[i]
			i--
		} else {
			nums1[k] = nums2[j]
			j--
		}
	}
}
