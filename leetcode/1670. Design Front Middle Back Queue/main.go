/*
Design a queue that supports push and pop operations in the front, middle, and back.

Implement the FrontMiddleBack class:

FrontMiddleBack() Initializes the queue.
void pushFront(int val) Adds val to the front of the queue.
void pushMiddle(int val) Adds val to the middle of the queue.
void pushBack(int val) Adds val to the back of the queue.
int popFront() Removes the front element of the queue and returns it. If the queue is empty, return -1.
int popMiddle() Removes the middle element of the queue and returns it. If the queue is empty, return -1.
int popBack() Removes the back element of the queue and returns it. If the queue is empty, return -1.

Notice that when there are two middle position choices,
the operation is performed on the frontmost middle position choice.

For example:
	Pushing 6 into the middle of [1, 2, 3, 4, 5] results in [1, 2, 6, 3, 4, 5].
	Popping the middle from [1, 2, 3, 4, 5, 6] returns 3 and results in [1, 2, 4, 5, 6].


Example 1:
	Input:
	["FrontMiddleBackQueue", "pushFront", "pushBack", "pushMiddle", "pushMiddle", "popFront", "popMiddle", "popMiddle", "popBack", "popFront"]
	[[], [1], [2], [3], [4], [], [], [], [], []]

	Output:
		[null, null, null, null, null, 1, 3, 4, 2, -1]

	Explanation:
		FrontMiddleBackQueue q = new FrontMiddleBackQueue();
		q.pushFront(1);   // [1]
		q.pushBack(2);    // [1, 2]
		q.pushMiddle(3);  // [1, 3, 2]
		q.pushMiddle(4);  // [1, 4, 3, 2]
		q.popFront();     // return 1 -> [4, 3, 2]
		q.popMiddle();    // return 3 -> [4, 2]
		q.popMiddle();    // return 4 -> [2]
		q.popBack();      // return 2 -> []
		q.popFront();     // return -1 -> [] (The queue is empty)

Constraints:
	1 <= val <= 109
	At most 1000 calls will be made to pushFront, pushMiddle, pushBack, popFront, popMiddle, and popBack.

*/

package main

type FrontMiddleBackQueue struct {
	left  []int
	right []int
}

func Constructor() FrontMiddleBackQueue {
	return FrontMiddleBackQueue{
		left:  []int{},
		right: []int{},
	}
}

func (q *FrontMiddleBackQueue) balance() {
	if len(q.left) > len(q.right)+1 {
		val := q.left[len(q.left)-1]
		q.left = q.left[:len(q.left)-1]
		q.right = append([]int{val}, q.right...)
	}

	if len(q.left) < len(q.right) {
		val := q.right[0]
		q.right = q.right[1:]
		q.left = append(q.left, val)
	}
}

func (q *FrontMiddleBackQueue) PushFront(val int) {
	q.left = append([]int{val}, q.left...)
	q.balance()
}

func (q *FrontMiddleBackQueue) PushMiddle(val int) {
	if len(q.left) > len(q.right) {
		lastLeft := q.left[len(q.left)-1]
		q.left = q.left[:len(q.left)-1]
		q.right = append([]int{lastLeft}, q.right...)
	}
	q.left = append(q.left, val)
}

func (q *FrontMiddleBackQueue) PushBack(val int) {
	q.right = append(q.right, val)
	q.balance()
}

func (q *FrontMiddleBackQueue) PopFront() int {
	if q.IsEmpty() {
		return -1
	}
	val := q.left[0]
	q.left = q.left[1:]
	q.balance()
	return val
}

func (q *FrontMiddleBackQueue) IsEmpty() bool {
	return len(q.left) == 0 && len(q.right) == 0
}

func (q *FrontMiddleBackQueue) PopMiddle() int {
	if q.IsEmpty() {
		return -1
	}

	val := q.left[len(q.left)-1]
	q.left = q.left[:len(q.left)-1]
	q.balance()
	return val
}

func (q *FrontMiddleBackQueue) PopBack() int {
	if q.IsEmpty() {
		return -1
	}
	var val int
	if len(q.right) == 0 {
		val = q.left[len(q.left)-1]
		q.left = q.left[:len(q.left)-1]
	} else {
		val = q.right[len(q.right)-1]
		q.right = q.right[:len(q.right)-1]
	}
	q.balance()
	return val
}

// type FrontMiddleBackQueue struct {
// 	queue []int
// }

// func Constructor() FrontMiddleBackQueue {
// 	return FrontMiddleBackQueue{
// 		queue: make([]int, 0),
// 	}
// }

// func (q *FrontMiddleBackQueue) PushFront(val int) {
// 	q.queue = append([]int{val}, q.queue...)
// }

// func (q *FrontMiddleBackQueue) PushMiddle(val int) {
// 	if q.IsEmpty() {
// 		q.queue = append(q.queue, val)
// 		return
// 	}
// 	mid := len(q.queue) / 2
// 	q.queue = append(q.queue[:mid], append([]int{val}, q.queue[mid:]...)...)
// }

// func (q *FrontMiddleBackQueue) PushBack(val int) {
// 	q.queue = append(q.queue, val)
// }

// func (q *FrontMiddleBackQueue) PopFront() int {
// 	if q.IsEmpty() {
// 		return -1
// 	}
// 	val := q.queue[0]
// 	q.queue = q.queue[1:]
// 	return val
// }

// func (q *FrontMiddleBackQueue) IsEmpty() bool {
// 	return len(q.queue) == 0
// }

// func (q *FrontMiddleBackQueue) PopMiddle() int {
// 	if q.IsEmpty() {
// 		return -1
// 	}
// 	mid := len(q.queue) / 2
// 	if len(q.queue)%2 == 0 {
// 		mid = mid - 1
// 	}
// 	val := q.queue[mid]
// 	if mid == 0 {
// 		q.queue = q.queue[1:]
// 	} else if mid == len(q.queue)-1 {
// 		q.queue = q.queue[:mid]
// 	} else {
// 		q.queue = append(q.queue[:mid], q.queue[mid+1:]...)
// 	}
// 	return val
// }

// func (q *FrontMiddleBackQueue) PopBack() int {
// 	if q.IsEmpty() {
// 		return -1
// 	}
// 	val := q.queue[len(q.queue)-1]
// 	q.queue = q.queue[:len(q.queue)-1]
// 	return val
// }
