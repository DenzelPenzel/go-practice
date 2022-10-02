/*
There are n workers. You are given two integer arrays quality and wage where quality[i] is the quality
of the ith worker and wage[i] is the minimum wage expectation for the ith worker.

We want to hire exactly k workers to form a paid group. To hire a group of k workers,
we must pay them according to the following rules:

Every worker in the paid group should be paid in the ratio of their quality compared to other workers in the paid group.
Every worker in the paid group must be paid at least their minimum wage expectation.
Given the integer k, return the least amount of money needed to form a paid group satisfying the above conditions.
Answers within 10-5 of the actual answer will be accepted.


Example 1:
    Input: quality = [10,20,5], wage = [70,50,30], k = 2
    Output: 105.00000
    Explanation: We pay 70 to 0th worker and 35 to 2nd worker.

Example 2:
    Input: quality = [3,1,10,10,1], wage = [4,8,2,2,7], k = 3
    Output: 30.66667
    Explanation: We pay 4 to 0th worker, 13.33333 to 2nd and 3rd workers separately.

Constraints:
    n == quality.length == wage.length
    1 <= k <= n <= 104
    1 <= quality[i], wage[i] <= 104
*/

package leetcode

import (
	"container/heap"
	"math"
	"sort"
)

type Worker struct {
	rate    float64
	quality int
	wage    int
}

type Item struct {
	quality int
}

type PQ []*Item

func (pq PQ) Len() int {
	return len(pq)
}

func (pq PQ) Less(i, j int) bool {
	return pq[i].quality > pq[j].quality
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PQ) Push(p interface{}) {
	*pq = append(*pq, p.(*Item))
}

func (pq *PQ) Pop() interface{} {
	p := (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return p
}

func mincostToHireWorkers(quality []int, wage []int, k int) float64 {
	n := len(quality)
	A := make([]*Worker, 0)
	res := math.MaxFloat64

	for i := 0; i < n; i++ {
		A = append(A, &Worker{rate: float64(wage[i]) / float64(quality[i]), quality: quality[i], wage: wage[i]})
	}

	sort.Slice(A, func(i, j int) bool {
		return A[i].rate < A[j].rate
	})

	running := 0

	queue := make(PQ, 0)

	for _, item := range A {
		running += item.quality

		heap.Push(&queue, &Item{quality: item.quality})

		for len(queue) > k {
			running -= heap.Pop(&queue).(*Item).quality
		}

		if len(queue) == k {
			res = math.Min(res, item.rate*float64(running))
		}
	}

	return float64(res)
}
