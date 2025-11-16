package main

import (
	"errors"
	"sync"
)

type RingBuffer struct {
	buffer   []interface{}
	capacity int
	head     int
	tail     int
	size     int
	mx       sync.Mutex
	cond     *sync.Cond
}

func NewRingBuffer(capacity int) *RingBuffer {
	if capacity < 0 {
		panic("capacity must be greater than 0")
	}
	rb := &RingBuffer{
		buffer:   make([]interface{}, capacity),
		capacity: capacity,
	}
	rb.cond = sync.NewCond(&rb.mx)
	return rb
}

func (rb *RingBuffer) Push(item interface{}) error {
	rb.mx.Lock()
	defer rb.mx.Unlock()

	if rb.size == rb.capacity {
		return errors.New("buffer is full")
	}

	rb.buffer[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.size++
	// Signal that an item has been added
	rb.cond.Signal()
	return nil
}

func (rb *RingBuffer) Pop() (interface{}, error) {
	rb.mx.Lock()
	defer rb.mx.Unlock()

	for rb.size == 0 {
		// Wait for an item to be available
		rb.cond.Wait()
	}

	v := rb.buffer[rb.head]
	// Zero-out the slot to avoid retaining references unnecessarily
	var zero interface{}
	rb.buffer[rb.head] = zero

	rb.head = (rb.head + 1) % rb.capacity
	rb.size--
	return v, nil
}

func (rb *RingBuffer) Size() int {
	rb.mx.Lock()
	defer rb.mx.Unlock()
	return rb.size
}

func (rb *RingBuffer) IsEmpty() bool {
	rb.mx.Lock()
	defer rb.mx.Unlock()
	return rb.size == 0
}

func (rb *RingBuffer) IsFull() bool {
	rb.mx.Lock()
	defer rb.mx.Unlock()
	return rb.size == rb.capacity
}

func (rb *RingBuffer) Clear() {
	rb.mx.Lock()
	defer rb.mx.Unlock()
	rb.head = 0
	rb.tail = 0
	rb.size = 0
}
