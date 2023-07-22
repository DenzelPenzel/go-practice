package main

import "sort"

type Item struct {
	x     int
	dir   string
	index int
}

func findRightInterval(A [][]int) []int {
	intervals := make([]Item, 0)

	for i, p := range A {
		intervals = append(intervals, Item{x: p[0], dir: "start", index: i})
		intervals = append(intervals, Item{x: p[1], dir: "end", index: i})
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].x != intervals[j].x {
			return intervals[i].x < intervals[j].x
		}

		return intervals[i].dir < intervals[j].dir
	})

	stack := make([]Item, 0)
	res := make([]int, len(A))

	for i := 0; i < len(A); i++ {
		res[i] = -1
	}

	for _, p := range intervals {
		if p.dir == "start" {
			for len(stack) > 0 && stack[len(stack)-1].x <= p.x {
				res[stack[0].index] = p.index
				stack = stack[1:]
			}
		} else {
			stack = append(stack, p)
		}
	}

	return res
}
