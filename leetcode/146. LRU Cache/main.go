/*
Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.

Implement the LRUCache class:

LRUCache(int capacity) Initialize the LRU cache with positive size capacity.
int get(int key) Return the value of the key if the key exists, otherwise return -1.
void put(int key, int value) Update the value of the key if the key exists. Otherwise, add the key-value pair to the cache.
If the number of keys exceeds the capacity from this operation, evict the least recently used key.
The functions get and put must each run in O(1) average time complexity.

Example 1:

Input
["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
[[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
Output
[null, null, null, 1, null, -1, null, -1, 3, 4]

Explanation
LRUCache lRUCache = new LRUCache(2);
lRUCache.put(1, 1); // cache is {1=1}
lRUCache.put(2, 2); // cache is {1=1, 2=2}
lRUCache.get(1);    // return 1
lRUCache.put(3, 3); // LRU key was 2, evicts key 2, cache is {1=1, 3=3}
lRUCache.get(2);    // returns -1 (not found)
lRUCache.put(4, 4); // LRU key was 1, evicts key 1, cache is {4=4, 3=3}
lRUCache.get(1);    // return -1 (not found)
lRUCache.get(3);    // return 3
lRUCache.get(4);    // return 4

Constraints:

1 <= capacity <= 3000
0 <= key <= 104
0 <= value <= 105
At most 2 * 105 calls will be made to get and put.
*/
package main

import "fmt"

type LRUNode struct {
	next, prev *LRUNode
	key, val   int
}

type LRUCache struct {
	head, tail *LRUNode
	mapping    map[int]*LRUNode
	capacity   int
}

func Constructor(capacity int) LRUCache {
	// create a pointer to the struct.
	head := &LRUNode{}
	tail := &LRUNode{}
	head.next = tail
	tail.prev = head

	return LRUCache{
		capacity: capacity,
		mapping:  map[int]*LRUNode{},
		head:     head,
		tail:     tail,
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

	// remove the last used node from double ll
	if len(lc.mapping) == lc.capacity {
		last_node := lc.last()

		node := lc.mapping[last_node.key]

		delete(lc.mapping, last_node.key)
		lc.remove(node)
	}

	newNode := new(LRUNode)
	newNode.key = key
	newNode.val = value
	lc.mapping[key] = newNode
	lc.insertFirst(newNode)
}

func (lc *LRUCache) remove(node *LRUNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (lc *LRUCache) last() *LRUNode {
	return lc.tail.prev
}

// insertFirst - insert a new node in the head of the double ll
func (lc *LRUCache) insertFirst(node *LRUNode) {
	node.next = lc.head.next
	node.prev = lc.head

	lc.head.next.prev = node
	lc.head.next = node
}

func main() {
	lr := Constructor(2)

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
}
