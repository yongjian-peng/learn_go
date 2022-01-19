package week01

import "fmt"

func RunReverseNodeInKGroup() {
	node := CreateNodeList([]int{1, 2, 3, 4, 5})

	ant := ReverseNodeInKGroup(node, 2)

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

func ReverseNodeInKGroup(head *ListNode, k int) *ListNode {
	if head == nil || k <= 1 {
		return head
	}

	var guard = new(ListNode)
	guard.Next = head
	var prev = guard

	for {
		var end = getEnd(head, k)
		if end == nil {
			break
		}
		var next = end.Next
		ReverseNode(head, k)
		prev.Next = end
		head.Next = next
		prev = head
		head = next
		if head == nil {
			break
		}
	}
	return guard.Next
}

func ReverseNode(head *ListNode, k int) {
	var prev = head
	head = head.Next
	k--
	for {
		var next = head.Next
		head.Next = prev
		k--
		if k == 0 {
			break
		}
		prev = head
		head = next
	}
	return
}

func getEnd(node *ListNode, k int) *ListNode {
	for {
		k--
		if k == 0 {
			return node
		}
		if node.Next == nil {
			break
		}
		node = node.Next
	}
	return nil
}
