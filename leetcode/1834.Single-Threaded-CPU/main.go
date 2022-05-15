package leetcode

import (
	"container/heap"
	"sort"
)

func getOrder(tasks [][]int) []int {
	t := make([]*Task, 0)

	for i := 0; i < len(tasks); i++ {
		t = append(t, &Task{tasks[i][0], tasks[i][1], i})
	}

	sort.Slice(t, func(i, j int) bool {
		if t[i].EnqueueTime == t[j].EnqueueTime {
			return t[i].ProcessTime < t[j].ProcessTime
		}
		return t[i].EnqueueTime < t[j].EnqueueTime

	})

	idx := 1
	time := t[0].EnqueueTime
	queue := make(PQ, 0)
	res := make([]int, 0)
	heap.Push(&queue, t[0])

	for idx < len(t) {
		cur := heap.Pop(&queue).(*Task)
		time += cur.ProcessTime
		res = append(res, cur.Index)

		for idx < len(t) && time >= t[idx].EnqueueTime {
			heap.Push(&queue, t[idx])
			idx += 1
		}

		if len(queue) == 0 && idx < len(t) {
			heap.Push(&queue, t[idx])
			time = t[idx].EnqueueTime
			idx += 1
		}
	}

	for len(queue) > 0 {
		cur := heap.Pop(&queue).(*Task)
		res = append(res, cur.Index)
	}

	return res
}

type Task struct {
	EnqueueTime int
	ProcessTime int
	Index       int
}

func max(i, j int) int {
	if i > j {
		return i
	}

	return j
}

type PQ []*Task

func (this PQ) Len() int {
	return len(this)
}

func (this PQ) Less(i, j int) bool {
	if this[i].ProcessTime == this[j].ProcessTime {
		return this[i].Index < this[j].Index
	} else {
		return this[i].ProcessTime < this[j].ProcessTime
	}
}

func (this PQ) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this *PQ) Push(p interface{}) {
	*this = append(*this, p.(*Task))
}

func (this *PQ) Pop() interface{} {
	p := (*this)[len(*this)-1]
	*this = (*this)[:len(*this)-1]

	return p

}
