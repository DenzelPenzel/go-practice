package main

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func countNodes(root *TreeNode) int {
	if root == nil {
		return 0
	}

	left, right := root, root
	d := 0

	for right != nil {
		left, right = left.Left, right.Right
		d++
	}

	if left == nil {
		return (1 << d) - 1
	}

	return 1 + countNodes(root.Left) + countNodes(root.Right)
}
