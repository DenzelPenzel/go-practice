package main

import "container/heap"

type node struct {
	u int
	w int
}

type Graph struct {
	graph map[int][]node
	n     int
}

type item struct {
	value    int
	priority int
	index    int
}

type priorityQueue []*item

func (pq *priorityQueue) Len() int { return len(*pq) }

func (pq *priorityQueue) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority
}

func (pq *priorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func Constructor(n int, edges [][]int) Graph {
	g := Graph{
		graph: make(map[int][]node),
	}
	for i := 0; i < n; i++ {
		g.graph[i] = []node{}
	}
	for _, edge := range edges {
		g.AddEdge(edge)
	}
	return g
}

func (g *Graph) AddEdge(edge []int) {
	v, u, w := edge[0], edge[1], edge[2]
	g.graph[v] = append(g.graph[v], node{u: u, w: w})
}

func (g *Graph) ShortestPath(node1 int, node2 int) int {
	dist := make(map[int]int)

	for k := range g.graph {
		dist[k] = 1<<31 - 1
	}

	dist[node1] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{value: node1, priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*item)

		if current.value == node2 {
			return current.priority
		}

		for _, neighbor := range g.graph[current.value] {
			newDist := current.priority + neighbor.w
			if newDist < dist[neighbor.u] {
				dist[neighbor.u] = newDist
				heap.Push(pq, &item{value: neighbor.u, priority: newDist})
			}
		}

	}

	return -1
}

/**
 * Your Graph object will be instantiated and called as such:
 * obj := Constructor(n, edges);
 * obj.AddEdge(edge);
 * param_2 := obj.ShortestPath(node1,node2);
 */
