package main

import (
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
func (q *LockFreeQueue[T]) Dequeue() {

}
