package utils

import (
	"fmt"
	f "fmt"
	fm "fmt"
	"math"
	"runtime"
	"sort"
	"strconv"
)

func minAvailableDuration(A [][]int, B [][]int, duration int) []int {
	// sort.Sort(sort.IntSlice(A))
	// sort.Sort(sort.IntSlice(B))
	res := make([][]int, 0)
	i, j := 0, 0

	for i < len(A) && j < len(B) {
		if overlap(A[i], B[j]) || overlap(B[j], A[i]) {
			x := max(A[i][0], B[j][0])
			y := min(A[i][1], B[j][1])
			res = append(res, []int{x, y})
		}
		if A[i][1] > B[j][1] {
			j++
		} else {
			i++
		}
	}
	// for  := range res {

	return []int{-1}
}

func Overlap(A []int, B []int) bool {
	return B[1] >= A[0] && A[1] >= B[0]
}

// func Min(a int, b int) int {
// 	if a > b {
// 		return b
// 	}
// 	return a
// }

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func neighbors(node string) []string {
	res := make([]string, 0)
	for i := 0; i < 4; i++ {
		x := int(node[i] - '0')
		for _, d := range []int{-1, 1} {
			y := (x + d) % 10
			res = append(res, node[:i]+strconv.Itoa(y)+node[i+1:])
		}
	}
	return res
}

func toposort(n int, relations [][]int) {
	graph := [][]int{}
	indegrees := []int{}

	for i := 0; i < n; i++ {
		graph = append(graph, []int{})
		indegrees = append(indegrees, 0)
	}

	for i := 0; i < len(relations); i++ {
		p := relations[i]
		graph[p[0]-1] = append(graph[p[0]-1], p[1]-1)
		indegrees[p[1]-1] += 1
	}

	queue := []int{}

	for i := 0; i < n; i++ {
		if indegrees[i] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		size := len(queue)

		for i := 0; i < size; i++ {
			v := queue[0]
			queue = queue[1:]
			for _, u := range graph[v] {
				indegrees[u] -= 1
				if indegrees[u] == 0 {
					queue = append(queue, u)
				}
			}
		}
	}
}

func minimumSemesters(n int, relations [][]int) int {
	graph := [][]int{}
	indegrees := []int{}

	for i := 0; i < n; i++ {
		graph = append(graph, []int{})
		indegrees = append(indegrees, 0)
	}

	for i := 0; i < len(relations); i++ {
		p := relations[i]
		graph[p[0]-1] = append(graph[p[0]-1], p[1]-1)
		indegrees[p[1]-1] += 1
	}

	queue := []int{}
	path := []int{}
	move := 0

	for i := 0; i < n; i++ {
		if indegrees[i] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			v := queue[0]
			queue = queue[1:]
			path = append(path, v)
			for _, u := range graph[v] {
				indegrees[u] -= 1
				if indegrees[u] == 0 {
					queue = append(queue, u)
				}
			}
		}
		move += 1
	}

	if len(path) != n {
		return -1
	}

	return move
}

func openLock(deadends []string, target string) int {
	if target == "0000" {
		return -1
	}

	visited := make(map[string]bool)
	dead := make(map[string]bool)
	for _, v := range deadends {
		visited[v] = true
		dead[v] = true
	}

	queue := []string{"0000"}
	move := 0

	for len(queue) > 0 {
		size := len(queue)

		for k := 0; k < size; k++ {
			node := queue[0]
			queue = queue[1:]

			if node == target {
				return move
			}

			if dead[node] {
				continue
			}

			for _, nei := range neighbors(node) {
				if !visited[nei] {
					visited[nei] = true
					queue = append(queue, nei)
				}
			}
		}
	}
	return -1
}

func queueTest(A [][]int) {
	n, m := len(A), len(A[0])
	queue := make([][]int, 0)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 0 {
				queue = append(queue, []int{i, j})
			}
		}
	}

	step := 0

	for len(queue) > 0 {
		s := len(queue)

		for k := 0; k < s; k++ {
			i, j := queue[0][0], queue[0][1]
			queue = queue[1:]

			for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
				x, y := p[0], p[1]

				if x >= 0 && x < n && y >= 0 && y < m {
					queue = append(queue, []int{x, y})
				}
			}

		}
		step += 1
	}

}

func main() {
	A := []int{3, 5, 1, 2, 10, -10}

	sort.Sort(sort.IntSlice(A))

	fmt.Println("Hello, World!", A, math.Inf(1), math.Max(10, 15))

	dict := map[string]int{"foo": 1, "bar": 2}
	value, ok := dict["baz"]
	if ok {
		fmt.Println("value: ", value)
	} else {
		fmt.Println("key not found")
	}

	fmt.Println("hello")
	f.Println("hey")
	fm.Println("hi")
}

func Min(a int, rest ...int) int {
	minValue := a
	for _, v := range rest {
		if v < minValue {
			minValue = v
		}
	}
	return minValue
}

func minimumDeletions(A []int) int {
	n := len(A)
	minval, maxval := math.Inf(1), math.Inf(-1)
	i, j := 0, 0

	for k, v := range v {
		if minval > v {
			minval = v
			i = k
		}
		if maxval < v {
			maxval = v
			j = k
		}
	}

	if i > j {
		i, j = j, i
	}

	removeFrontBack := i + 1 + n - j
	removeFront := i + 1 + j - i
	removeBack := n - j + j - i

	return Min(removeFrontBack, removeFront, removeBack)
}

func Version() string {
	return runtime.Version()
}
