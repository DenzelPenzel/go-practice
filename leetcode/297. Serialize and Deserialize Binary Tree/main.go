/*
Serialization is the process of converting a data structure or object into a sequence of bits
so that it can be stored in a file or memory buffer, or transmitted across a network connection link
to be reconstructed later in the same or another computer environment.

Design an algorithm to serialize and deserialize a binary tree.
There is no restriction on how your serialization/deserialization algorithm should work.
You just need to ensure that a binary tree can be serialized to a string and
this string can be deserialized to the original tree structure.

Clarification:
    The input/output format is the same as how LeetCode serializes a binary tree.
    You do not necessarily need to follow this format,
    so please be creative and come up with different approaches yourself.

Example 1:
  Input: root = [1,2,3,null,null,4,5]
  Output: [1,2,3,null,null,4,5]

Example 2:
  Input: root = []
  Output: []

Example 3:
  Input: root = [1]
  Output: [1]

Example 4:
  Input: root = [1,2]
  Output: [1,2]

Constraints:
  The number of nodes in the tree is in the range [0, 104].
  -1000 <= Node.val <= 1000
*/

package leetcode

import (
	"strconv"
	"strings"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type Codec struct {
	builder strings.Builder
	input   []string
}

func Constructor() Codec {
	return Codec{}
}

// Serializes a tree to a single string.
func (this *Codec) serialize(root *TreeNode) string {
	if root == nil {
		this.builder.WriteString("# ")
		return ""
	}
	this.builder.WriteString(strconv.Itoa(root.Val) + " ")
	this.serialize(root.Left)
	this.serialize(root.Right)
	return this.builder.String()
}

// Deserializes your encoded data to tree.
func (this *Codec) deserialize(data string) *TreeNode {
	if len(data) == 0 {
		return nil
	}
	this.input = strings.Split(data, " ")
	return this.DFS()
}

func (this *Codec) DFS() *TreeNode {
	if this.input[0] == "#" {
		this.input = this.input[1:]
		return nil
	}
	val, _ := strconv.Atoi(this.input[0])
	this.input = this.input[1:]
	return &TreeNode{
		Val:   val,
		Left:  this.DFS(),
		Right: this.DFS(),
	}
}
