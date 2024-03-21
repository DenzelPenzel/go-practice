package main

// Bitwise Trie
type BitwiseTrie struct {
	children [2]*BitwiseTrie
	val      int
}

func NewBitwiseTrie() *BitwiseTrie {
	return &BitwiseTrie{children: [2]*BitwiseTrie{nil, nil}}
}

func findMaximumXOR(A []int) int {
	res := 0

	trie := NewBitwiseTrie()

	for _, x := range A {
		node := trie

		for i := 31; i >= 0; i-- {
			bit := (x >> i) & 1
			if node.children[bit] == nil {
				node.children[bit] = NewBitwiseTrie()
			}
			node = node.children[bit]
		}
		// save number in the node
		node.val = x
	}

	for _, x := range A {
		node := trie

		for i := 31; i >= 0; i-- {
			bit := (x >> i) & 1
			rev := bit ^ 1

			if node.children[rev] != nil {
				node = node.children[rev]
			} else {
				node = node.children[bit]
			}
		}

		res = max(res, x^node.val)
	}

	return res
}
