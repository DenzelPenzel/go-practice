package main

import (
	"fmt"
	"math"
)

type state struct {
	vis  map[int]bool
	path []int
	dist map[int]int
	step int
}

func longestCycle(edges []int) int {
	n := len(edges)
	mapping := make(map[int]int, n)
	ans := math.MinInt32

	ss := &state{
		vis:  make(map[int]bool, n),
		path: []int{},
		dist: make(map[int]int, n),
		step: 0,
	}

	for v, u := range edges {
		mapping[v] = u
	}

	for v := 0; v < n; v++ {
		if _, ok := ss.vis[v]; !ok {
			ss.step = 0
			ss.dfs(v, mapping)
			dec := 1

			for len(ss.path) > 0 {
				vv := ss.path[0]

				if len(ss.path) > 1 && vv == ss.path[len(ss.path)-1] {
					dec = 0
					ans = max(ans, ss.step)
				}

				ss.dist[vv] = ss.step
				ss.step -= dec
				ss.path = ss.path[1:]
			}
		}
	}

	fmt.Println(ss.dist)

	if ans == math.MinInt32 {
		return -1
	}

	return ans
}

func (s *state) dfs(v int, mapping map[int]int) {
	s.path = append(s.path, v)
	if _, ok := s.vis[v]; ok {
		s.step += s.dist[v]
		return
	}
	s.vis[v] = true
	s.step += 1
	if u, ok := mapping[v]; ok {
		s.dfs(u, mapping)
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
