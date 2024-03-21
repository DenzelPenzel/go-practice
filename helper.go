package main

import "github.com/emirpasic/gods/maps/treemap"

// @formatter:off
var mod = int(1e9) + 7
var dirs = [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var length = int(^uint(0) >> 1) // Set length to max int value

// Clearing maps
func clearMaps(m map[string]int) {
	for k := range m {
		delete(m, k)
	}
}

// Fill the array 
A := make([]int, 10)
tmp := A
fmt.Println(A) // [0,0,0,0,0,0,0,0,0,0]
for _, c := range [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}} {
	copy(tmp, c) // [1,2,3,0,0,0,0,0,0,0]
	tmp = tmp[len(c):] // [0,0,0,0,0,0,0]
}
fmt.Println(nums)


// Increasing the length of a slice
func incSlice() {
	s := make([]int, 100)
	if len(s) < 200 {
		// Need to extend s
		s = append(s, make([]int, 100)...)
	}
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func absI(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func cloneIntArray(old []int) []int { clone := make([]int, len(old)); copy(clone, old); return clone }

func createInt1dSlice(params ...int) []int {
	n := params[0]
	d := 0
	if len(params) == 2 {
		d = params[1]
	}
	_ret := make([]int, n)
	for i := range _ret {
		_ret[i] = d
	}
	return _ret
}
func createInt2dSlice(params ...int) [][]int {
	m := params[0]
	n := params[1]
	d := 0
	if len(params) == 3 {
		d = params[2]
	}
	_ret := make([][]int, m)
	for i := range _ret {
		_ret[i] = createInt1dSlice(n, d)
	}
	return _ret
}
func createBool2dSlice(m int, n int) [][]bool {
	_ret := make([][]bool, m)
	for i := range _ret {
		_ret[i] = make([]bool, n)
	}
	return _ret
}
func getSieve(n int) []bool {
	_sieve := make([]bool, n+1)
	_sieve[0], _sieve[1] = true, true
	for p := 2; p*p <= n; p++ {
		if !_sieve[p] {
			for i := p * p; i <= n; i += p {
				_sieve[i] = true
			}
		}
	}
	return _sieve
}
func lcm(a int, b int) int { return (a * b) / gcd(a, b) }
func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
func nCr(n, r int) int {
	prod := 1
	r = minI(r, n-r)
	for i := 1; i <= r; i++ {
		prod *= n - i + 1
		prod /= i
	}
	return prod
}
func signI(a int) int {
	if a < 0 {
		return -1
	} else if a == 0 {
		return 0
	} else {
		return 1
	}
}
func powMod(a, b, mod int64) int64 {
	result := int64(1)
	a = a % mod
	for b > 0 {
		if b&1 == 1 {
			result = (result * a) % mod
		}
		a = (a * a) % mod
		b >>= 1
	}
	return result
}
func pow(a, b int) int {
	result := 1
	for b > 0 {
		if b&1 == 1 {
			result *= a
		}
		a *= a
		b >>= 1
	}
	return result
}

type Bit struct{ tree []int }

func NewBit(size int) *Bit { return &Bit{tree: make([]int, size+1)} }
func (b *Bit) Sum(index int) int {
	sum := 0
	index++
	for index > 0 {
		sum += b.tree[index]
		index -= index & -index
	}
	return sum
}
func (b *Bit) Update(index, value int) {
	index++
	for index < len(b.tree) {
		b.tree[index] += value
		index += index & -index
	}
}

type UnionFind struct {
	parent     []int
	rank       []int
	components int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	components := n
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank, components: components}
}
func (uf *UnionFind) Find(i int) int {
	if uf.parent[i] != i {
		uf.parent[i] = uf.Find(uf.parent[i])
	}
	return uf.parent[i]
}
func (uf *UnionFind) Union(i, j int) bool {
	p1 := uf.Find(i)
	p2 := uf.Find(j)
	if p1 == p2 {
		return false
	}
	if uf.rank[p1] < uf.rank[p2] {
		uf.parent[p1] = p2
	} else if uf.rank[p2] < uf.rank[p1] {
		uf.parent[p2] = p1
	} else {
		uf.parent[p2] = p1
		uf.rank[p1]++
	}
	uf.components--
	return true
}

type Seg struct {
	tree []int
	n    int
}

func NewSeg(n int) *Seg { return &Seg{tree: make([]int, 2*n), n: n} }
func NewSegFromArray(array []int) *Seg {
	n := len(array)
	tree := make([]int, n*2)
	copy(tree[n:], array)
	for i := n - 1; i > 0; i-- {
		tree[i] = tree[i<<1] + tree[i<<1|1]
	}
	return &Seg{tree: tree, n: n}
}
func (s *Seg) Sum(l, r int) int {
	l, r, sum := l+s.n, r+s.n, 0
	for l <= r {
		if l&1 == 1 {
			sum += s.tree[l]
			l++
		}
		if r&1 == 0 {
			sum += s.tree[r]
			r--
		}
		l >>= 1
		r >>= 1
	}
	return sum
}
func (s *Seg) Update(index, value int) {
	index += s.n
	s.tree[index] += value
	for ; index > 1; index >>= 1 {
		s.tree[index>>1] = s.tree[index] + s.tree[index^1]
	}
}
func incrementKey(tree *treemap.Map, key int) {
	value, found := tree.Get(key)
	if !found {
		value = 0
	}
	tree.Put(key, value.(int)+1)
}
func decrementKey(tree *treemap.Map, key int) {
	value, found := tree.Get(key)
	if found {
		value = value.(int) - 1
		if value == 0 {
			tree.Remove(key)
		} else {
			tree.Put(key, value)
		}
	}
}

type Stack []int

func NewStack(capacity ...int) *Stack {
	c := 0
	if len(capacity) == 1 {
		c = capacity[0]
	}
	var stack Stack = make([]int, 0, c)
	return &stack
}
func (s *Stack) Push(num int) { *s = append(*s, num) }
func (s *Stack) Pop() int     { val := s.Top(); *s = (*s)[:len(*s)-1]; return val }
func (s *Stack) Top() int     { return (*s)[len(*s)-1] }
func (s *Stack) Empty() bool  { return s.Size() == 0 }
func (s *Stack) Size() int    { return len(*s) }

// @formatter:off
