
type pair struct {
	node *TreeNode
	path string
}

func getDirections(root *TreeNode, startValue int, destValue int) string {
	lca := getLCA(root, startValue, destValue)

	st := make([]pair, 0)
	st = append(st, pair{lca, ""})

	destPath := ""
	startDepth := 0

	for len(st) > 0 {
		p := st[len(st)-1]
		st = st[:len(st)-1]

		if p.node.Val == startValue {
			startDepth = len(p.path)
		}

		if p.node.Val == destValue {
			destPath = p.path
		}

		if startDepth > 0 && len(destPath) > 0 {
			break
		}

		if p.node.Left != nil {
			st = append(st, pair{p.node.Left, p.path + "L"})
		}

		if p.node.Right != nil {
			st = append(st, pair{p.node.Right, p.path + "R"})
		}
	}

	start := ""
	if startDepth > 0 {
		for i := 0; i < startDepth; i++ {
			start += "U"
		}
	}
	return start + destPath
}

func getLCA(node *TreeNode, p int, q int) *TreeNode {
	if node == nil || node.Val == p || node.Val == q {
		return node
	}

	left := getLCA(node.Left, p, q)
	right := getLCA(node.Right, p, q)

	if left != nil && right != nil {
		return node
	}

	if left != nil {
		return left
	}

	return right
}
