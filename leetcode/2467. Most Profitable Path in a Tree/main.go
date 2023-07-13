/*
2467. Most Profitable Path in a Tree (Medium)

There is an undirected tree with n nodes labeled from 0 to n - 1, rooted at
node 0. You are given a 2D integer array edges of length n - 1 where
edges[i] = [ai, bi] indicates that there is an edge between nodes ai and bi
in the tree. At every node i, there is a gate. You are also given an array
of even integers amount, where amount[i] represents:
* the price needed to open the gate at node i, if amount[i] is negative, or,
* the cash reward obtained on opening the gate at node i, otherwise.

The game goes on as follows:
* Initially, Alice is at node 0 and Bob is at node bob.
* At every second, Alice and Bob each move to an adjacent node. Alice moves
    towards some leaf node, while Bob moves towards node 0.
* For every node along their path, Alice and Bob either spend money to open
    the gate at that node, or accept the reward. Note that:
    + If the gate is already open, no price will be required, nor will there
    be any cash reward.
    + If Alice and Bob reach the node simultaneously, they share the
    price/reward for opening the gate there. In other words, if the price
    to open the gate is c, then both Alice and Bob pay c / 2 each.
    Similarly, if the reward at the gate is c, both of them receive c / 2
    each.
* If Alice reaches a leaf node, she stops moving. Similarly, if Bob reaches
    node 0, he stops moving. Note that these events are independent of each
    other.

Return the maximum net income Alice can have if she travels towards the
optimal leaf node.

Example 1:
Input: edges = [[0,1],[1,2],[1,3],[3,4]], bob = 3, amount = [-2,4,2,-4,6]
Output: 6
Explanation: The above diagram represents the given tree. The game goes as
                follows:
                - Alice is initially on node 0, Bob on node 3. They open the
                gates of their respective nodes. Alice's net income is now
                -2.
                - Both Alice and Bob move to node 1. Since they reach here
                simultaneously, they open the gate together and share the
                reward. Alice's net income becomes -2 + (4 / 2) = 0.
                - Alice moves on to node 3. Since Bob already opened its gate,
                Alice's income remains unchanged. Bob moves on to node 0,
                and stops moving.
                - Alice moves on to node 4 and opens the gate there. Her net
                income becomes 0 + 6 = 6.
                Now, neither Alice nor Bob can make any further moves, and the
                game ends. It is not possible for Alice to get a higher net
                income.

Example 2:
Input: edges = [[0,1]], bob = 1, amount = [-7280,2350]
Output: -7280
Explanation: Alice follows the path 0->1 whereas Bob follows the path 1->0.
                Thus, Alice opens the gate at node 0 only. Hence, her net
                income is -7280.

Constraints:
* 2 <= n <= 10^5
* edges.length == n - 1
* edges[i].length == 2
* 0 <= ai, bi < n
* ai != bi
* edges represents a valid tree.
* 1 <= bob < n
* amount.length == n
* amount[i] is an even integer in the range [-10^4, 10^4].
*/

package main

import "math"

func mostProfitablePath(edges [][]int, bob int, amount []int) int {
	n := len(edges) + 1
	graph := make([][]int, n)
	parent := make([]int, n)
	aliceDist := make([]int, n)

	for _, p := range edges {
		graph[p[0]] = append(graph[p[0]], p[1])
		graph[p[1]] = append(graph[p[1]], p[0])
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
		res := math.MinInt32
		for _, u := range graph[v] {
			if u != p {
				res = maxI(res, amount[u]+dfs2(u, v))
			}
		}
		return res
	}

	dfs(0, -1, 0)

	v := bob
	bobDist := 0

	for v != 0 {
		if aliceDist[v] == bobDist {
			amount[v] = amount[v] / 2
		} else if aliceDist[v] > bobDist {
			amount[v] = 0
		}
		v = parent[v]
		bobDist++
	}

	return amount[0] + dfs2(0, -1)
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
