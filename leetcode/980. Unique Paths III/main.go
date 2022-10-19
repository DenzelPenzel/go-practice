package leetcode

func uniquePathsIII(A [][]int) int {
	n, m := len(A), len(A[0])
	var start_x int
	var start_y int
	empty := 0

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 1 {
				start_x, start_y = i, j
			}
			if A[i][j] == 0 {
				empty += 1
			}
		}
	}

	return DFS(&A, start_x, start_y, empty)
}

func DFS(A *[][]int, i, j, empty int) int {
	if (*A)[i][j] == 2 {
		if empty == -1 {
			return 1
		}
		return 0
	}

	res := 0
	(*A)[i][j] = -1

	for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
		x, y := p[0], p[1]
		if x >= 0 && x < len(*A) && y >= 0 && y < len((*A)[0]) && (*A)[x][y] != -1 {
			res += DFS(A, x, y, empty-1)
		}
	}

	(*A)[i][j] = 0

	return res
}
