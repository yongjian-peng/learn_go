// 在这个示例中，我们定义了 SortStrategy 接口和两个具体的排序策略，BubbleSortStrategy 和 QuickSortStrategy。然后，我们定义了 Context 结构体来包含当前的排序策略，并提供了 SetStrategy 和 Sort 方法来设置和执行排序策略。
//
//在 main 函数中，我们创建了一个 Context 实例，并分别使用冒泡排序和快速排序策略来对数组进行排序。最后，我们使用 Go 内置的 sort.Ints 函数对数组进行排序，以便进行比较。

package main

import (
	"fmt"
	"sort"
)

// SortStrategy 是排序策略接口
type SortStrategy interface {
	Sort([]int) []int
}

// BubbleSortStrategy 是冒泡排序策略
type BubbleSortStrategy struct{}

// Sort 实现 BubbleSortStrategy 接口
func (s BubbleSortStrategy) Sort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	return arr
}

// QuickSortStrategy 是快速排序策略
type QuickSortStrategy struct{}

// Sort 实现 QuickSortStrategy 接口
func (s QuickSortStrategy) Sort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	pivot := arr[0]
	left, right := []int{}, []int{}
	for _, v := range arr[1:] {
		if v <= pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	left = s.Sort(left)
	right = s.Sort(right)
	return append(append(left, pivot), right...)
}

// Context 是排序上下文，用于根据不同的排序策略执行不同的排序算法
type Context struct {
	strategy SortStrategy
}

// NewContext 创建排序上下文
func NewContext(strategy SortStrategy) *Context {
	return &Context{strategy: strategy}
}

// SetStrategy 设置排序策略
func (c *Context) SetStrategy(strategy SortStrategy) {
	c.strategy = strategy
}

// Sort 对数组进行排序
func (c *Context) Sort(arr []int) []int {
	return c.strategy.Sort(arr)
}

func main() {
	arr := []int{5, 2, 8, 1, 9, 4}
	context := NewContext(BubbleSortStrategy{})
	fmt.Println(context.Sort(arr))

	context.SetStrategy(QuickSortStrategy{})
	fmt.Println(context.Sort(arr))

	sort.Ints(arr)
	fmt.Println(arr)
}
