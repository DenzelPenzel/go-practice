package leetcode

import "sort"

func findOriginalArray(changed []int) []int {
	if len(changed) <= 0 || len(changed)%2 != 0 {
		return []int{}
	}

	count := make(map[int]int)
	res := []int{}

	for _, v := range changed {
		count[v] += 1
	}

	sort.Ints(changed)

	for _, x := range changed {
		if count[x] == 0 {
			continue
		}
		y := x * 2
		if v, ok := count[y]; ok && v > 0 {
			res = append(res, x)
			count[x] -= 1
			count[y] -= 1
		} else {
			return []int{}
		}
	}

	return res

}
