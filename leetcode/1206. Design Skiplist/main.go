package leetcode

import (
	"math"
	"math/rand"
)

type ListNode struct {
	val, cnt   int
	next, down *ListNode
}

type Skiplist struct {
	head *ListNode
	prob float32
}

func Constructor() Skiplist {
	return Skiplist{
		head: &ListNode{
			val:  math.MinInt64,
			cnt:  1,
			next: nil,
			down: nil,
		},
		prob: 0.25,
	}
}

func (this *Skiplist) Search(target int) bool {
	node := this.head
	for node != nil && node.val < target {
		if node.next != nil && node.next.val == target {
			return true
		}
		if node.next != nil && node.next.val < target {
			node = node.next
		} else {
			node = node.down
		}
	}
	return false
}

func (this *Skiplist) Add(num int) {
	node := this.head
	st := make([]*ListNode, 0)

	for node != nil && node.val < num {
		if node.next != nil && node.next.val <= num {
			node = node.next
		} else {
			st = append(st, node)
			node = node.down
		}
	}

	if node != nil {
		for node != nil {
			node.cnt += 1
			node = node.down
		}
	} else {
		var prev *ListNode

		for len(st) != 0 {
			if len(st) != 0 {
				node = st[len(st)-1]
				newNode := &ListNode{val: num, cnt: 1, next: node.next, down: prev}
				node.next = newNode
				prev = newNode
				st = st[:len(st)-1]
			} else {
				newNode := &ListNode{val: math.MinInt64, cnt: 1, next: nil, down: this.head}
				this.head = newNode
				nextNode := &ListNode{val: num, down: prev, cnt: 1, next: nil}
				this.head.next = prev
				prev = nextNode
			}

			r := rand.Float32()

			if r >= this.prob {
				break
			}
		}
	}
}

func (this *Skiplist) Erase(num int) bool {
	node := this.head
	st := make([]*ListNode, 0)

	for node != nil {
		if node.next != nil && node.next.val < num {
			node = node.next
		} else {
			st = append(st, node)
			node = node.down
		}
	}

	res := false

	for len(st) > 0 {
		node := st[len(st)-1]

		if node.next != nil && node.next.val == num {
			res = true
			if node.next.cnt > 1 {
				node.next.cnt -= 1
			} else {
				node.next = node.next.next
			}
		} else {
			break
		}
	}

	return res
}
