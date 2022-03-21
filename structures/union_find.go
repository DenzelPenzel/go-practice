package structures

type UnionFind struct {
	parent, rank []int
	count        int
}

func (this *UnionFind) Init(count int) {
	this.count = count
	this.parent = make([]int, count)
	this.rank = make([]int, count)

	for i := range this.parent {
		this.parent[i] = i
	}
}

func (this *UnionFind) Find(p int) int {
	root := p

	for root != this.parent[root] {
		root = this.parent[root]
	}
	// compress path
	for p != this.parent[p] {
		tmp := this.parent[p]
		this.parent[p] = root
		p = tmp
	}
	return root
}

// Union define
func (this *UnionFind) Union(p, q int) {
	proot := this.Find(p)
	qroot := this.Find(q)
	if proot == qroot {
		return
	}
	if this.rank[qroot] > this.rank[proot] {
		this.parent[proot] = qroot
	} else {
		this.parent[qroot] = proot
		if this.rank[proot] == this.rank[qroot] {
			this.rank[proot]++
		}
	}
	this.count--
}

// TotalCount define
func (this *UnionFind) TotalCount() int {
	return this.count
}

// UnionFindCount define
type UnionFindCount struct {
	parent, count []int
	maxUnionCount int
}

func (this *UnionFindCount) Init(n int) {
	this.parent = make([]int, n)
	this.count = make([]int, n)
	for i := range this.parent {
		this.parent[i] = i
		this.count[i] = i
	}
}

// Find define
func (uf *UnionFindCount) Find(p int) int {
	root := p
	for root != uf.parent[root] {
		root = uf.parent[root]
	}
	return root
}

// Union define
func (uf *UnionFindCount) Union(p, q int) {
	proot := uf.Find(p)
	qroot := uf.Find(q)
	if proot == qroot {
		return
	}
	if proot == len(uf.parent)-1 {
		//proot is root
	} else if qroot == len(uf.parent)-1 {
		// qroot is root, always attach to root
		proot, qroot = qroot, proot
	} else if uf.count[qroot] > uf.count[proot] {
		proot, qroot = qroot, proot
	}

	//set relation[0] as parent
	uf.maxUnionCount = max(uf.maxUnionCount, (uf.count[proot] + uf.count[qroot]))
	uf.parent[qroot] = proot
	uf.count[proot] += uf.count[qroot]
}

// Count define
func (uf *UnionFindCount) Count() []int {
	return uf.count
}

// MaxUnionCount define
func (uf *UnionFindCount) MaxUnionCount() int {
	return uf.maxUnionCount
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
