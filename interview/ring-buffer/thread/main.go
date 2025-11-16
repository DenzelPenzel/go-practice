package main

type RingBuffer struct {
	capacity int
}

func NewRingBuffer(capacity int) *RingBuffer {
	if capacity < 0 {
		panic("capacity must be greater than 0")
	}
	rb := &RingBuffer{
		capacity: capacity,
	}
	return rb
}

func main() {
}
