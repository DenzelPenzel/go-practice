package main

import "container/heap"

type Node struct {
	value int
	i, j  int
}

type PQ []*Node

func (pq PQ) Len() int {
	return len(pq)
}

func (pq PQ) Less(i, j int) bool {
	return pq[i].value < pq[j].value
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PQ) Push(p interface{}) {
	*pq = append(*pq, p.(*Node))
}

func (pq *PQ) Pop() interface{} {
	p := (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return p
}

type state struct {
	i, j int
}

func kSmallestPairs(nums1 []int, nums2 []int, k int) [][]int {
	n, m := len(nums1), len(nums2)
	queue := make(PQ, 0)
	vis := make(map[state]bool, n)
	heap.Push(&queue, &Node{value: nums1[0] + nums2[0], i: 0, j: 0})
	var res [][]int

	for k > 0 && len(queue) > 0 {
		node := heap.Pop(&queue).(*Node)
		res = append(res, []int{nums1[node.i], nums2[node.j]})

		if node.i == n && node.j == m {
			break
		}

		dirs := [][]int{[]int{0, 1}, []int{1, 0}}

		for _, dir := range dirs {
			next_i, next_j := node.i+dir[0], node.j+dir[1]
			s := state{i: next_i, j: next_j}

			if next_i < n && next_j < m && !vis[s] {
				vis[s] = true
				heap.Push(&queue, &Node{value: nums1[next_i] + nums2[next_j], i: next_i, j: next_j})
			}
		}
		k--
	}

	return res
}
