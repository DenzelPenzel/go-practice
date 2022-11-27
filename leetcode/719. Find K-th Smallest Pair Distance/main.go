package main

import "sort"

func smallestDistancePair(A []int, k int) int {
	sort.Sort(sort.IntSlice(A))

	lo, hi := 0, (A[len(A)-1]-A[0])+1

	for lo < hi {
		dist := lo + ((hi - lo) >> 1)
		cntPairs := getCountPairs(A, dist)

		if cntPairs >= k {
			hi = dist
		} else {
			lo = dist + 1
		}
	}

	return lo
}

func getCountPairs(A []int, target int) int {
	res := 0
	j := 0
	for i := 0; i < len(A); i++ {
		for j < len(A) && A[j]-A[i] <= target {
			j++
		}
		res += j - i - 1
	}
	return res
}
