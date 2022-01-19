package week01

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func CreateNodeList(nums []int) *ListNode {
	if len(nums) == 0 {
		return nil
	}
	var guard = &ListNode{}

	var cur = guard
	for i := 0; i < len(nums); i++ {
		var node = &ListNode{Val: nums[i]}
		cur.Next = node
		cur = node
	}
	return guard.Next
}

func RunReverseList() {
	list := CreateNodeList([]int{1, 2, 3, 4, 5})

	ant := ReverseList(list)

	for {
		if ant != nil {
			fmt.Println(ant.Val)
			ant = ant.Next
		} else {
			break
		}
	}

	fmt.Println("OK")
}

func ReverseList(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	var cur = head
	var prev *ListNode
	for {
		next := cur.Next
		cur.Next = prev
		if next == nil {
			break
		}
		prev = cur
		cur = next
	}
	return cur
}
