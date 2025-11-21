package main

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"time"
)

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
				res = max(res, amount[u]+dfs2(u, v))
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

func main() {
	// logger := logging.New(time.RFC3339, true)

	// logger.Log("info", "starting up service")
	// logger.Log("warning", "no tasks found")
	// logger.Log("error", "exiting: no work performed")

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
		fmt.Println("--", v)
	}

	// slice descriptor contains a pointer to an underlying array, along with length and capacity
	// pointer to a slice (&aaa) -> references the slice descriptor, not the underlying array
	aaa := []int{7, 8, 9}
	for i, v := range *(&aaa) {
		fmt.Println(i, v)
	}

	array := [6]int{0, 1, 2, 3, 4, 5}
	slice := array[1:3] // [1 2]

	fmt.Println(cap(slice), len(slice), slice)

	//

	var v int = 10
	fmt.Println(v)

	//
	ch := make(chan int)

	go func() {
		time.Sleep(time.Second * 2)
		close(ch)
	}()

	fmt.Println("Waiting for channel to be closed...")
	<-ch
	fmt.Println("Channel closed")
}
