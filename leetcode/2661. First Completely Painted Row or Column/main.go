/*
You are given a 0-indexed integer array arr, and an m x n integer matrix mat. arr and mat both contain all the integers
in the range [1, m * n].

Go through each index i in arr starting from index 0 and paint the cell in mat containing the integer arr[i].

Return the smallest index i at which either a row or a column will be completely painted in mat.

Example 1:

	image explanation for example 1
	Input: arr = [1,3,4,2], mat = [[1,4],[2,3]]
	Output: 2
	Explanation: The moves are shown in order, and both the first row and second column of the matrix
	become fully painted at arr[2].

Example 2:

	image explanation for example 2
	Input: arr = [2,8,7,4,1,3,5,6,9], mat = [[3,2,5],[1,4,6],[8,7,9]]
	Output: 3
	Explanation: The second column becomes fully painted at arr[3].

Constraints:

	m == mat.length
	n = mat[i].length
	arr.length == m * n
	1 <= m, n <= 105
	1 <= m * n <= 105
	1 <= arr[i], mat[r][c] <= m * n
	All the integers of arr are unique.
	All the integers of mat are unique.
*/
package main

func firstCompleteIndex(A []int, mat [][]int) int {
	n, m := len(mat), len(mat[0])
	mp := make(map[int][2]int)
	rows := make([]int, n)
	cols := make([]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			mp[mat[i][j]] = [2]int{i, j}
		}
	}

	for k, val := range A {
		p := mp[val]
		rows[p[0]] += 1
		cols[p[1]] += 1
		if rows[p[0]] == m || cols[p[1]] == n {
			return k
		}
	}

	return -1
}
