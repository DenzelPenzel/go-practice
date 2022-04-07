package leetcode

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
