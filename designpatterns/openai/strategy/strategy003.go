package main

import "fmt"

// 在这个示例中，我们定义了 SortStrategy3 接口和两个具体的排序策略， BubbleSortStrategy3 和 QuickSortStrategy3 。然后，我们定义了 Context3 结构体来包含当前的排序策略，并提供了 SetStrategy3 和 Sort3 方法来设置和执行排序策略

// 在 main 函数中，我们创建了一个 Context 实例，并分别使用冒泡排序和快速排序策略来对数组进行排序，
// 最后 使用 Go 内置的 sort.Ints 函数对数组进行排序，以便进行比较

// SortStrategy3 是排序策略接口
type SortStrategy3 interface {
	Sort3([]int) []int
}

// BubbleSortStrategy3 是冒泡排序策略
type BubbleSortStrategy3 struct{}

// Sort 实现 BubbleSortStrategy3 接口
func (s BubbleSortStrategy3) Sort3(arr []int) []int {
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

// QuickSortStrategy3 是快速排序策略
type QuickSortStrategy3 struct{}

// Sort3 实现 QuickSortStrategy3 接口
func (s QuickSortStrategy3) Sort3(arr []int) []int {
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
	left = s.Sort3(left)
	right = s.Sort3(right)
	return append(append(left, pivot), right...)
}

// Context3 是排序上下文，用于根据不同的排序策略执行不同的排序算法
type Context3 struct {
	strategy3 SortStrategy3
}

// NewContext3 创建排序上下文
func NewContext3(strategy3 SortStrategy3) *Context3 {
	return &Context3{strategy3: strategy3}
}

// SetStrategy3 设置排序策略
func (c *Context3) SetStrategy3(strategy3 SortStrategy3) {
	c.strategy3 = strategy3
}

// Sort3 对数组进行排序
func (c *Context3) Sort3(arr []int) []int {
	return c.strategy3.Sort3(arr)
}

func main() {
	arr := []int{5, 2, 8, 1, 9, 4}

	context := NewContext3(BubbleSortStrategy3{})
	fmt.Println(context.Sort3(arr))

	context.SetStrategy3(QuickSortStrategy3{})
	fmt.Println(context.Sort3(arr))
}
