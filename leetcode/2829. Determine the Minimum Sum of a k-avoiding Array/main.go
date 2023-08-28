package main

func minimumSum(n int, k int) int {
	dfs := func(start int, comb map[int]bool, sum int, *res int) {
		if len(comb) == n {
			*res = maxI(*res, sum)
			return
		}

		if sum > *res {
			return
		}

		if _, ok := comb[k - start]; !ok {
			comb[start] = true
			dfs(start + 1, comb, sum + start, &res)
			delete(comb[start])
		} else {
			dfs(start + 1, comb, sum, &res)
		}
	}

	res := 0

	dfs(1, []int{}, 0, &res)


	return res
}


func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
