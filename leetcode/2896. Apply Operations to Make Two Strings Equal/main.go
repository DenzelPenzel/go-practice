package main

func minOperations(s1 string, s2 string, x int) int {
	var dfs func([]int, int, float32) float32
	dfs = func(diffs []int, i int, float32) float32 {
		if i == 0 {
			return x / 2
		}
		if i == -1 {
			return 0
		}

		removeBothCost := diffs[i] - diffs[i-1]
		removeLastCost := x / 2

		return minI(dfs(diffs, i-2, x)+removeBothCost, dfs(diffs, i-1, x)+removeLastCost)
	}

	diffs := make([]int, 0)

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			diffs = append(diffs, i)
		}
	}

	if len(diffs)%2 != 0 {
		return -1
	}

	res := dfs(diffs, len(diffs)-1, float32(x))

	return int(res)
}

func minI(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
