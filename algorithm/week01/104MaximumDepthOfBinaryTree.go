package week01

import (
	"github.com/halfrost/LeetCode-Go/structures"
)

type TreeNode = structures.TreeNode

func RunMaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(RunMaxDepth(root.Left), RunMaxDepth(root.Right)) + 1
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
