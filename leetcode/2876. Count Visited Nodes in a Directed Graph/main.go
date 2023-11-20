/*
There is a directed graph consisting of n nodes numbered from 0 to n - 1 and n directed edges.

You are given a 0-indexed array edges where edges[i] indicates that there is an edge from node i to node edges[i].

Consider the following process on the graph:

You start from a node x and keep visiting other nodes through edges until you reach a node that you have already visited before on this same process.
Return an array answer where answer[i] is the number of different nodes that you will visit if you perform the process starting from node i.

Example 1:
	Input: edges = [1,2,0,0]
	Output: [3,3,3,4]
	Explanation: We perform the process starting from each node in the following way:
	- Starting from node 0, we visit the nodes 0 -> 1 -> 2 -> 0. The number of different nodes we visit is 3.
	- Starting from node 1, we visit the nodes 1 -> 2 -> 0 -> 1. The number of different nodes we visit is 3.
	- Starting from node 2, we visit the nodes 2 -> 0 -> 1 -> 2. The number of different nodes we visit is 3.
	- Starting from node 3, we visit the nodes 3 -> 0 -> 1 -> 2 -> 0. The number of different nodes we visit is 4.

Example 2:
	Input: edges = [1,2,3,4,0]
	Output: [5,5,5,5,5]
	Explanation: Starting from any node we can visit every node in the graph in the process.


Constraints:
	n == edges.length
	2 <= n <= 10^5
	0 <= edges[i] <= n - 1
	edges[i] != i
*/

package main

func countVisitedNodes(edges []int) []int {
	n := len(edges)
	dp := make([]int, n)
	vis := make(map[int]bool, 0)
	var path []int

	var dfs func(int, int) int
	dfs = func(v int, steps int) int {
		path = append(path, v)
		if _, ok := vis[v]; ok {
			steps += dp[v]
			return steps
		}
		vis[v] = true
		return dfs(edges[v], steps+1)
	}

	for v := 0; v < n; v++ {
		if _, ok := vis[v]; !ok {
			steps := dfs(v, 0)
			dec := 1
			for len(path) > 0 {
				// find the cycle, the degree is the same for all vertices
				if path[0] == path[len(path)-1] {
					dec = 0
				}
				vv := path[0]
				path = path[1:len(path)]
				dp[vv] = steps
				steps -= dec
			}
		}
	}

	return dp

}
