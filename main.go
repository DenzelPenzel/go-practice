package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
)

func maximumProfit(present []int, future []int, budget int) int {
	A := make([][]int, 0)
	for i := 0; i < len(future); i++ {
		if future[i]-present[i] > 0 {
			A = append(A, []int{present[i], future[i] - present[i]})
		}
	}

	if len(A) == 0 {
		return 0
	}

	n := len(A)
	dp := [][]int{}

	for i := 0; i < n; i++ {
		dp = append(dp, []int{})
		for j := 0; j <= budget; j++ {
			dp[i] = append(dp[i], 0)
		}
	}

	for w := 0; w <= budget; w++ {
		if A[0][0] <= w {
			dp[0][w] = A[0][1]
		}
	}

	for i := 1; i < n; i++ {
		for w := 0; w <= budget; w++ {
			if A[i][0] <= w {
				dp[i][w] = Max(dp[i-1][w], dp[i-1][w-A[i][0]]+A[i][1])
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	return dp[len(dp)-1][budget]
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

type Shop struct {
	dist  int
	price int
	row   int
	col   int
}

func highestRankedKItems(grid [][]int, pricing []int, start []int, k int) [][]int {
	n, m := len(grid), len(grid[0])
	seen := make([][]bool, n)
	queue := make([][]int, 0)
	queue = append(queue, []int{start[0], start[1]})
	items := make([]Shop, 0)
	dist := 0

	for i := range seen {
		seen[i] = make([]bool, m)
	}

	for len(queue) > 0 {
		size := len(queue)

		for k := 0; k < size; k++ {
			i, j := queue[0][0], queue[0][1]
			queue = queue[1:]

			if pricing[0] <= grid[i][j] && grid[i][j] <= pricing[1] {
				items = append(items, Shop{dist, grid[i][j], i, j})
			}

			for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
				x, y := p[0], p[1]

				if x >= 0 && x < n && y >= 0 && y < m {
					if grid[x][y] >= 1 {
						queue = append(queue, []int{x, y})
						seen[x][y] = true
					}
				}
			}
		}
		dist++
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].dist != items[j].dist {
			return items[i].dist < items[j].dist
		}
		if items[i].price != items[j].price {
			return items[i].price < items[j].price
		}
		if items[i].row != items[j].row {
			return items[i].row < items[j].row
		}
		return items[i].col < items[j].col
	})

	res := make([][]int, 0, k)
	for i := 0; i < len(items) && i < k; i++ {
		res = append(res, []int{items[i].row, items[i].col})
	}
	return res
}

func shortestBridge(A [][]int) int {
	n, m := len(A), len(A[0])
	stack := make([][]int, 0)
	seen := make([][]bool, n)

	for i := range seen {
		seen[i] = make([]bool, m)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 1 {
				stack = append(stack, []int{i, j})
				break
			}
		}
		if len(stack) > 0 {
			break
		}
	}

	for len(stack) > 0 {
		i, j := stack[len(stack)-1][0], stack[len(stack)-1][1]
		stack = stack[:len(stack)-1]
		for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j + 1}, {i, j - 1}} {
			x, y := p[0], p[1]

			if 0 <= x && x < n && 0 <= y && y < m && A[x][y] == 1 && !seen[x][y] {
				seen[x][y] = true
				stack = append(stack, []int{x, y})
			}
		}
	}

	queue := make([][]int, 0)
	for i, v := range seen {
		for j := range v {
			if v[j] {
				queue = append(queue, []int{i, j})
			}
		}
	}

	step := 0

	for len(queue) > 0 {
		size := len(queue)

		for k := 0; k < size; k++ {
			i, j := queue[0][0], queue[0][1]
			queue := queue[1:]

			for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j + 1}, {i, j - 1}} {
				x, y := p[0], p[1]

				if 0 <= x && x < n && 0 <= y && y < m && !seen[x][y] {
					if A[x][y] == 1 {
						return step
					}
					seen[x][y] = true
					queue = append(queue, []int{x, y})
				}
			}
		}
		step += 1
	}

	return step
}

type ListNode struct {
	val, cnt   int
	next, down *ListNode
}

type Skiplist struct {
	head *ListNode
	prob float32
}

func CConstructor() Skiplist {
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

type LRUNode struct {
	next, prev *LRUNode
	key, val   int
}

type LRUCache struct {
	head, tail *LRUNode
	mapping    map[int]*LRUNode
	capacity   int
}

func Constructor(capacity int) LRUCache {
	head := new(LRUNode)
	tail := new(LRUNode)
	head.next = tail
	tail.prev = head

	return LRUCache{
		capacity: capacity,
		mapping:  map[int]*LRUNode{},
		head:     head,
		tail:     tail,
	}
}

func (lc *LRUCache) Get(key int) int {
	if node, ok := lc.mapping[key]; ok {
		lc.remove(node)
		lc.append(node)
		return node.val
	}

	return -1
}

func (lc *LRUCache) Put(key int, value int) {
	if lc.capacity == 0 {
		return
	}

	if node, ok := lc.mapping[key]; ok {
		node.val = value

		lc.remove(node)
		lc.append(node)
		return
	}

	// remove the last used node from double ll
	if len(lc.mapping) == lc.capacity {
		last_node := lc.last()

		node := lc.mapping[last_node.key]

		delete(lc.mapping, last_node.key)
		lc.remove(node)
	}

	newNode := new(LRUNode)
	lc.mapping[key] = newNode
	lc.append(newNode)
}

func (lc *LRUCache) remove(node *LRUNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (lc *LRUCache) last() *LRUNode {
	return lc.tail.prev
}

// append - insert a new node in the head of the double ll
func (lc *LRUCache) append(node *LRUNode) {
	node.next = lc.head.next
	node.prev = lc.head

	lc.head.next.prev = node
	lc.head.next = node
}

func test() {
	type user struct {
		name     string
		username string
	}

	var users map[string]user

	// panic
	users["vasya"] = user{name: "vasya", username: "vas"}
}

func inspectSlice(slice []string) {
	fmt.Println("=======")
	fmt.Printf("Length[%d] Capacity[%d]\n", len(slice), cap(slice))
	for i := range slice {
		fmt.Printf("[%d] %p %s\n", i, &slice[i], slice[i])
	}
}

func mostProfitablePath(edges [][]int, bob int, amount []int) int {
	n := len(edges) + 1
	graph := make([][]int, n)
	parent := make([]int, n)
	aliceDist := make([]int, n)

	for _, p := range edges {
		graph[p[0]] = append(graph[p[0]], p[1])
	}

	var dfs func(v, p, d int)
	dfs = func(v, p, d int) {
		parent[v] = p
		aliceDist[v] = d
		for _, u := range graph[v] {
			if u != p {
				dfs(u, v, d+1)
			}
		}
	}

	var dfs2 func(v, p int) int
	dfs2 = func(v, p int) int {
		if len(graph[v]) == 1 && graph[v][0] == p {
			return 0
		}
		res := -math.MinInt32
		for _, u := range graph[v] {
			if u != p {
				res = maxI(res, amount[u]+dfs2(u, v))
			}
		}
		return res
	}

	dfs(0, -1, 0)

	v := bob
	d := 0

	for v != 0 {
		if aliceDist[v] == d {
			amount[v] = amount[v] / 2
		} else {
			amount[v] = 0
		}
		v = parent[v]
		d++
	}

	return amount[0] + dfs2(0, -1)
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func shortestBeautifulSubstring(s string, k int) string {
	var res []string
	length := int(^uint(0) >> 1)
	ones := 0
	i := 0
	for j, x := range s {
		if x == '1' {
			ones++
		}
		for ones >= k {
			if length > j-i+1 {
				length = j - i + 1
				res = []string{s[i : j+1]}
			} else if length == j-i+1 {
				res = append(res, s[i:j+1])
			}
			if s[i] == '1' {
				ones--
			}
			i++
		}
	}

	sort.Strings(res)

	if len(res) > 0 {
		return ""
	}

	return res[0]
}

func lcm(a int, b int) int {
	return (a * b) / gcd(a, b)
}

func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func main() {
	// logger := logging.New(time.RFC3339, true)

	// logger.Log("info", "starting up service")
	// logger.Log("warning", "no tasks found")
	// logger.Log("error", "exiting: no work performed")

	// fmt.Println(maximumProfit([]int{5, 4, 6, 2, 3}, []int{8, 5, 4, 3, 5}, 10))

	// fmt.Println(leetcode.GetOrder([][]int{{1, 2}, {2, 4}, {3, 2}, {4, 1}}))

	slist := CConstructor()
	slist.Add(1)
	slist.Add(2)
	slist.Add(3)
	fmt.Println("val: %t", slist.Search(0))
	slist.Add(4)
	fmt.Println("val: %t", slist.Search(1))
	fmt.Println("val: %t", slist.Erase(0))
	fmt.Println("val: %t", slist.Erase(1))
	fmt.Println("val: %t", slist.Search(1))

	A := []int{1, 3, 2}
	B := make([]int, len(A))
	copy(B, A)
	sort.Ints(B)
	same := reflect.ValueOf(A).Pointer() == reflect.ValueOf(B).Pointer()
	fmt.Println(same)
	fmt.Println(A, B)

	aa := make([]string, 0)
	aa = append(aa, "a")
	aa = append(aa, "b")
	aa = append(aa, "c")
	aa = append(aa, "d")
	inspectSlice(aa) // len: 4 cap: 4

	// small capacity
	// in this case ```append``` creates a new backing array (doubling or growing by 25%)
	// and then copies the values from the `old` array into the `new` one
	aa = append(aa, "e") // len: 5 cap: 8
	inspectSlice(aa)

	// s := "世界"
	// fmt.Println(len(s), utf8.RuneLen(s))

	str := "Hello, 世界"
	for index, r := range str {
		fmt.Printf("Character: %c, Unicode: %U, Byte position: %d\n", r, r, index)
	}

	// Go makes a copy, but the copy is hidden from us, and only ```v``` can see that copied array
	a := [3]int{1, 2, 3}
	b := [3]int{4, 5, 6}
	for i, v := range a {
		if i == 1 {
			a = b
		}
		fmt.Println(v)
	}

	// slice descriptor contains a pointer to an underlying array, along with length and capacity
	// pointer to a slice (&aaa) -> references the slice descriptor, not the underlying array

	aaa := []int{7, 8, 9}
	for i, v := range *(&aaa) {
		fmt.Println(i, v)
	}

}
