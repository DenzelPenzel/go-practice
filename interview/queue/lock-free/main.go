package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type node[T any] struct {
	value T
	next  unsafe.Pointer
}

type LockFreeQueue[T any] struct {
	head unsafe.Pointer // *node[T]
	tail unsafe.Pointer // *node[T]
}

// The Michael-Scott algorithm is susceptible to the ABA problem,
// where a pointer's value changes from A to B and back to A, misleading concurrent operations
func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	dummy := unsafe.Pointer(&node[T]{})
	return &LockFreeQueue[T]{
		head: dummy,
		tail: dummy,
	}
}

// Enqueue add a new element to the end of the queue
func (q *LockFreeQueue[T]) Enqueue(value T) {
	newNode := &node[T]{value: value}
	newNodePtr := unsafe.Pointer(newNode)

	for {
		tail := (*node[T])(atomic.LoadPointer(&q.tail))
		next := atomic.LoadPointer(&tail.next)

		if tail == (*node[T])(atomic.LoadPointer(&q.tail)) {
			if next == nil {
				// Tail is pointing to the last node
				if atomic.CompareAndSwapPointer(&tail.next, next, newNodePtr) {
					// Enqueue done; try to move the tail
					atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), newNodePtr)
					return
				}
			} else {
				// Tail not pointing to the last node; try to move it forward
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), next)
			}
		}

	}
}

// Dequeue removes and returns a value from the front of the queue.
// The boolean return value indicates whether the dequeue was successful.
func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	var zero T

	for {
		head := (*node[T])(atomic.LoadPointer(&q.head))
		tail := (*node[T])(atomic.LoadPointer(&q.tail))
		next := (*node[T])(atomic.LoadPointer(&head.next))

		if head == (*node[T])(atomic.LoadPointer(&q.head)) {
			if head == nil {
				// Queue is empty
				return zero, false
			}
			// Tail is falling behind; try to move it forward
			atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
		} else {
			// Read the value before CAS, as another dequeue might free the node
			value := next.value
			// Try to move head forward
			if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
				return value, true
			}
		}
	}
}

func main() {
	q := NewLockFreeQueue[int]()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			q.Enqueue(val)
		}(i)
	}

	wg.Wait()

	// Dequeue elements
	count := 0
	for {
		val, ok := q.Dequeue()
		if !ok {
			break
		}
		fmt.Println(val)
		count++
	}

	fmt.Printf("Total dequeued: %d\n", count)
}
