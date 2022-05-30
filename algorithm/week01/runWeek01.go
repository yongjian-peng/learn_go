package week01

import (
	"fmt"

	"github.com/halfrost/LeetCode-Go/structures"
)

func RunWeek01() {
	// 柱状图中最大矩形
	// RunLargestRectangleArea()
	// 合并两个有序数组
	// RunMergeTwoArray()

	// 滑动窗口最大值
	// RunMaxSlidingWindow()

	// 接雨水
	// RunTrap()

	// 反转链表
	// RunReverseList()

	// k 个一组反转链表
	// RunReverseNodeInKGroup()

	// 选择排序
	// RunSelectSort()

	// 二叉树的最大深度
	para104 := []int{3, 9, 20, structures.NULL, structures.NULL, 15, 7}
	root := structures.Ints2TreeNode(para104)
	depth := RunMaxDepth(root)

	fmt.Println(depth)

}
