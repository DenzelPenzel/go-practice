package main

import "fmt"

type Node struct {
	next, prev *Node
	key, val   int
}

type LRUCache struct {
	head, tail *Node
	mapping    map[int]*Node
	capacity   int
}

func NewLRUCache(capacity int) *LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return &LRUCache{
		head:     head,
		tail:     tail,
		mapping:  make(map[int]*Node),
		capacity: capacity,
	}
}

func (lc *LRUCache) Get(key int) int {
	if node, ok := lc.mapping[key]; ok {
		lc.remove(node)
		lc.insertFirst(node)
		return node.val
	}
	return -1
}

func (lc *LRUCache) Put(key int, value int) {
	if lc.capacity == 0 {
		return
	}

	if node, ok := lc.mapping[key]; ok {
		node.val = value
		lc.remove(node)
		lc.insertFirst(node)
		return
	}

	if len(lc.mapping) == lc.capacity {
		last_node := lc.last()

		node := lc.mapping[last_node.key]

		delete(lc.mapping, last_node.key)
		lc.remove(node)
	}

	newNode := &Node{key: key, val: value}
	lc.mapping[key] = newNode
	lc.insertFirst(newNode)
}

func (lc *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (lc *LRUCache) last() *Node {
	return lc.tail.prev
}

func (lc *LRUCache) insertFirst(node *Node) {
	node.next = lc.head.next
	node.prev = lc.head

	lc.head.next.prev = node
	lc.head.next = node
}

func main() {
	lr := NewLRUCache(2)

	lr.Put(1, 1)
	lr.Put(2, 2)
	x := lr.Get(1)
	fmt.Println("----", x)
	lr.Put(3, 3)
	lr.Get(3)
	lr.Get(2)
	lr.Put(4, 4)
	lr.Get(1)
	lr.Get(3)
	lr.Get(4)
	fmt.Println("----", lr.Get(4))
}
