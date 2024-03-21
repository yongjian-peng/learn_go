package main

import "fmt"

type ListNode struct {
	Value int
	Next  *ListNode
}

func reverseList(head *ListNode) *ListNode {
	// 定义空节点和当前节点
	var prev *ListNode
	curr := head

	// 遍历整个链表
	for curr != nil {
		// 保存当前节点
		tmp := curr.Next

		// 将当前节点的 Next 指向前一个节点
		curr.Next = prev

		// 更新 prev 和 curr
		prev = curr
		curr = tmp
	}
	return prev
}

// 翻转链表
func main() {
	// 创建一个链表
	l := &ListNode{Value: 1}

	l.Next = &ListNode{Value: 2}
	l.Next.Next = &ListNode{Value: 3}
	l.Next.Next.Next = &ListNode{Value: 4}

	for curr := l; curr != nil; curr = curr.Next {
		fmt.Println(curr.Value, " -> ")
	}

	newHead := reverseList(l)

	for curr := newHead; curr != nil; curr = curr.Next {
		fmt.Println(curr.Value, " -> ")
	}

}
