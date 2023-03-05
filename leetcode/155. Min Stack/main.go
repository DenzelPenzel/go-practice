/*
Design a stack that supports push, pop, top, and retrieving the minimum element in constant time.

Implement the MinStack class:

MinStack() initializes the stack object.
void push(int val) pushes the element val onto the stack.
void pop() removes the element on the top of the stack.
int top() gets the top element of the stack.
int getMin() retrieves the minimum element in the stack.
You must implement a solution with O(1) time complexity for each function.

Example 1:

Input
["MinStack","push","push","push","getMin","pop","top","getMin"]
[[],[-2],[0],[-3],[],[],[],[]]

Output
[null,null,null,null,-3,null,0,-2]

Explanation
MinStack minStack = new MinStack();
minStack.push(-2);
minStack.push(0);
minStack.push(-3);
minStack.getMin(); // return -3
minStack.pop();
minStack.top();    // return 0
minStack.getMin(); // return -2

Constraints:

-231 <= val <= 231 - 1
Methods pop, top and getMin operations will always be called on non-empty stacks.
At most 3 * 104 calls will be made to push, pop, top, and getMin.
*/
package main

type Item struct {
	val, cnt int
}

type MinStack struct {
	mapping []Item
	items   []int
}

func Constructor() MinStack {
	return MinStack{
		items:   []int{},
		mapping: []Item{},
	}
}

func (ms *MinStack) Push(val int) {
	if len(ms.mapping) <= 0 {
		ms.items = append(ms.items, val)
		ms.mapping = append(ms.mapping, Item{val: val, cnt: 1})
		return
	}

	if ms.mapping[len(ms.mapping)-1].val == val {
		ms.mapping[len(ms.mapping)-1].cnt += 1
	} else if len(ms.mapping) > 0 && ms.mapping[len(ms.mapping)-1].val > val {
		ms.mapping = append(ms.mapping, Item{val: val, cnt: 1})
	}

	ms.items = append(ms.items, val)
}

func (ms *MinStack) Pop() {
	val := ms.items[len(ms.items)-1]
	ms.items = ms.items[:len(ms.items)-1]
	if len(ms.mapping) > 0 && ms.mapping[len(ms.mapping)-1].val == val {
		ms.mapping[len(ms.mapping)-1].cnt -= 1
		// pop last elem if cnt == 0
		if ms.mapping[len(ms.mapping)-1].cnt <= 0 {
			ms.mapping = ms.mapping[:len(ms.mapping)-1]
		}
	}
}

func (ms *MinStack) Top() int {
	return ms.items[len(ms.items)-1]
}

func (ms *MinStack) GetMin() int {
	return ms.mapping[len(ms.mapping)-1].val
}
