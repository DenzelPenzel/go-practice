/*
You are given a 0-indexed integer array nums. A pair of integers x and y is called a strong pair
if it satisfies the condition:

|x - y| <= min(x, y)
You need to select two integers from nums such that they form a strong pair and
their bitwise XOR is the maximum among all strong pairs in the array.

Return the maximum XOR value out of all possible strong pairs in the array nums.

Note that you can pick the same integer twice to form a pair.

Example 1:
	Input: nums = [1,2,3,4,5]
	Output: 7
	Explanation:
		There are 11 strong pairs in the array nums: (1, 1), (1, 2), (2, 2), (2, 3), (2, 4),
		(3, 3), (3, 4), (3, 5), (4, 4), (4, 5) and (5, 5).
		The maximum XOR possible from these pairs is 3 XOR 4 = 7.

Example 2:
	Input: nums = [10,100]
	Output: 0
	Explanation:
		There are 2 strong pairs in the array nums: (10, 10) and (100, 100).
		The maximum XOR possible from these pairs is 10 XOR 10 = 0 since the
		pair (100, 100) also gives 100 XOR 100 = 0.

Example 3:
	Input: nums = [500,520,2500,3000]
	Output: 1020
	Explanation:
		There are 6 strong pairs in the array nums: (500, 500), (500, 520), (520, 520),
		(2500, 2500), (2500, 3000) and (3000, 3000).
		The maximum XOR possible from these pairs is 500 XOR 520 = 1020 since the only
		other non-zero XOR value is 2500 XOR 3000 = 636.

Constraints:
	1 <= nums.length <= 5 * 10^4
	1 <= nums[i] <= 220 - 1
*/

package main

import (
	"math"
	"sort"
)

type Node struct {
	children map[int]*Node
	val      int
	isEnd    bool
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{
		root: &Node{
			children: make(map[int]*Node),
			val:      0,
			isEnd:    false,
		},
	}
}

func (t *Trie) insert(num int) {
	node := t.root
	for i := 31; i >= 0; i-- {
		bit := (num >> i) & 1
		if _, ok := node.children[bit]; !ok {
			node.children[bit] = &Node{
				children: make(map[int]*Node),
				val:      0,
				isEnd:    false,
			}
		}
		node = node.children[bit]
	}
	node.val = num
	node.isEnd = true
}

func (t *Trie) search(num int) *Node {
	node := t.root
	for i := 31; i >= 0; i-- {
		bit := (num >> i) & 1
		rev := bit ^ 1
		if _, ok := node.children[rev]; ok {
			node = node.children[rev]
		} else {
			node = node.children[bit]
		}
	}
	return node
}

func (t *Trie) delete(num int) {
	var dfs func(node *Node, i int) bool
	dfs = func(node *Node, i int) bool {
		if i == -1 {
			if !node.isEnd {
				return false
			}
			node.isEnd = false
			return len(node.children) == 0
		}
		bit := (num >> i) & 1
		if node.children[bit] == nil {
			return false
		}
		shouldRemove := dfs(node.children[bit], i-1)
		if shouldRemove {
			delete(node.children, bit)
			return len(node.children) == 0
		}
		return false
	}
	dfs(t.root, 31)
}

func uniqueSorted(nums []int) []int {
	set := make(map[int]bool)
	for _, num := range nums {
		set[num] = true
	}

	var unique []int
	for num := range set {
		unique = append(unique, num)
	}

	sort.Ints(unique)
	return unique
}

func maximumStrongPairXor(nums []int) int {
	res := 0
	A := uniqueSorted(nums)

	trie := NewTrie()
	start, end := 0, 0

	for end < len(A) {
		for start < end && int(math.Abs(float64(A[start])-float64(A[end]))) > min(A[start], A[end]) {
			trie.delete(A[start])
			start++
		}
		trie.insert(A[end])
		if end-start > 0 {
			node := trie.search(A[end])
			if node != nil {
				res = max(res, A[end]^node.val)
			}
		}
		end++
	}

	return res
}
