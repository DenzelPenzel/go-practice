package leetcode

func findAllRecipes(recipes []string, ingredients [][]string, supplies []string) []string {
	graph := make(map[string][]string)
	indegrees := map[string]int{}
	n := len(recipes)
	seen := make(map[string]bool)

	for _, v := range recipes {
		seen[v] = true
	}

	for i := 0; i < n; i++ {
		for _, v := range ingredients[i] {
			graph[v] = append(graph[v], recipes[i])

			if indegrees[recipes[i]] > 0 {
				indegrees[recipes[i]] += 1
			} else {
				indegrees[recipes[i]] = 1
			}
		}
	}

	queue := make([]string, 0, len(recipes))
	res := []string{}

	for _, v := range supplies {
		queue = append(queue, v)
	}

	for len(queue) > 0 {
		size := len(queue)

		for k := 0; k < size; k++ {
			v := queue[0]
			queue = queue[1:]
			if _, ok := seen[v]; ok {
				res = append(res, v)
			}
			for _, u := range graph[v] {
				indegrees[u] -= 1
				if indegrees[u] == 0 {
					queue = append(queue, u)
				}
			}
		}
	}

	return res
}
