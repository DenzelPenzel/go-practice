package structures

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_intHeap(t *testing.T) {
	ast := assert.New(t)

	init_heap := new(intHeap)
	heap.Init(init_heap)

	heap.Push(init_heap, 1)
	heap.Pop(init_heap)

	begin, end := 0, 10

	for i := begin; i < end; i++ {
		heap.Push(init_heap, i)
		ast.Equal(0, (*init_heap)[0], "The minimum value after inserting %d is %d, init_heap=%v", i, (*init_heap)[0], (*init_heap))
	}

	for i := begin; i < end; i++ {
		fmt.Println(i, *init_heap)
		ast.Equal(i, heap.Pop(init_heap), "Pop init_heap=%v", (*init_heap))
	}

}
