## Go practice
The repository contains a collection of algorithmic and architectural questions related to the programming language Golang.


### Data types
- boolean 
- int 
  - signed and unsigned
- floating-point 
  - float32, float64
- Complex numbers 
  - complex64, complex128
- Byte 
  - alias for uint8
- Rune 
  - alias for int32, represents a Unicode code point
- String
  - represents a sequence of bytes
- Arrays
  - Fixed-size sequences of elements of a single type
- Slices
  - Dynamic, flexible view into the elements of an array
- Structs
    ```
    type Person struct {
        Name string
        Age  int
    }
    ```
- Maps
  - Unordered collections of key-value pairs ```map[string]int```

- Channels
  - communication between goroutines ```chan int```

- Pointers
  - Holds the memory address of another variable ```*int``` is a pointer to an integer

### Synchronization primitives
- Mutex
  - sync.Mutex
  - sync.RWMutex - multiple readers or one writer at a time

- WaitGroup
  
- Cond
  - used to block goroutines until a condition is met
    ```
        var mu sync.Mutex
        var cond = sync.NewCond(&mu)

        go func() {
            mu.Lock()
            cond.Wait() // Wait for a condition
            // work after condition is met
            mu.Unlock()
        }()

        mu.Lock()
        // prepare condition
        cond.Signal() // Wake one goroutine waiting on cond
        mu.Unlock()
    ```

- Once
  - code is executed only once, even if called from multiple goroutines
  ```
    var once sync.Once
    once.Do(func() {
        // initialize something
    })
  ```

- Atomic Operations (low-level atomic memory primitives)
    ```
    var counter int64
    atomic.AddInt64(&counter, 1) // Atomic increment
    atomic.LoadInt64(&counter) // Atomic read
    atomic.StoreInt64(&counter, 42) // Atomic write
    atomic.CompareAndSwapInt64(&counter, old, new) // CAS operation
    ```

### Zero Value Concept

Every single value I construct in Go is initialized at least to its zero value state
unless I specify the initialization value at construction.

The zero value is the setting of every bit in every byte to zero

### Padding and Alignment

How much memory is allocated for a value of type example?

```

type example struct {
    flag bool
    counter int16
    pi float32
}
```

bool - 1 bytes + int16 - 2 bytes + float 32 - 4 bytes = 7 bytes, however actual size is 8 bytes.

Why?
Bcs there is a padding byte sitting between the flag and counter fields for the reason of alignment.

Idea of alignment is to allow the hardware to read memory
more efficiently by placing memory on specific alignment boundaries.

### Padding example

```
type example2 struct {
    flag bool // 0xc000100020 <- Starting Address
    byte // 0xc000100021 <- 1 byte padding
    counter int16 // 0xc000100022 <- 2 byte alignment
    flag2 bool // 0xc000100024 <- 1 byte alignment
    byte // 0xc000100025 <- 1 byte padding
    byte // 0xc000100026 <- 1 byte padding
    byte // 0xc000100027 <- 1 byte padding
    pi float32 // 0xc000100028 <- 4 byte alignment
}

```

Solution: If I need to minimize the amount of padding bytes, I must lay out the fields from
highest allocation to lowest allocation

### Pointers

- Data in golang is moved by value

### Strings

- Immutable, strings are byte sequences, Go uses UTF-8 encoding by default
- Rune type is an alias for int32 used to represent a Unicode code point
- Unicode code point is a unique number assigned to each character 

```
// Convert string to rune slice
str = "Hello, 世界"
runes := []rune(str)
fmt.Println(runes) // Output: [72 101 108 108 111 44 32 19990 30028]
```

### Channels

- Send Operation: send a value into a non-buffered channel, it will block until another goroutine is ready to receive that value

- Receive Operation: receive from a non-buffered channel, it will block until another goroutine sends a value into that channel

```
package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	go func() {
        // operation blocks if there is no corresponding receive operation waiting
        // remain blocked until another goroutine executes <-ch to receive the value
		ch <- 42 
		fmt.Println("Sent value")
	}()

	value := <-ch // Receiving goroutine will block here until the sending goroutine sends
	fmt.Println("Received value:", value)
}
```


### Goroutine

- Goroutine occupies a few KB, this can support a large number of threads in a limited memory space goroutine
- Goroutine use less memory space and reduces the cost of context switching
- Each Goroutine is given its own block of memory called a stack
- Each stack starts out as a 2048 byte (2k) allocation
- Func is called 
  - Allocation of the stack space to execute func 
  - This block of memory is called a frame
- Size of a frame for a given function is calculated at compile time
- If the compiler doesn’t know the size of a value at compile time, the value has to be constructed on the heap
- Concept of coroutines allows a set of reusable functions to run on a set of threads
- Even if a coroutine blocks, other coroutines of the thread can be scheduled and transferred runtime to other runnable threads

### Problems with large number of threads

- High memory usage
- High CPU consumption for scheduling

### G scheduler

- it use G-M-P model
- G — represents goroutine, which is a task to be executed
- M — represents the thread of the operating system, which is scheduled and managed by the scheduler of the operating system
- P — represents the processor, which can be thought of as a local scheduler running on a thread
- GOMAXPROCS set to the number of cores of the current machine
- 4 active operating system threads will be created on a 4 core machine

### CPU caches

Prefetcher attempts to predict what data is needed before instructions request the data so it’s already present in either the L1 or L2 cache.

Program can read/write a byte of memory as the smallest unit of memory access.
Caching systems granularity is 64 bytes. 
This 64 byte block of memory is called a cache line.

Prefetcher works best when the instructions being executed create predictable access patterns to memory. One way to create a predictable access pattern to memory is to construct a contiguous block of memory and then iterate over that memory performing a linear traversal with a predictable stride.

Array is data structure with predictable access patterns. However, the slice is the most important data structure in Go. Slices in Go use an array underneath.

Prefetcher will pick up predictable data access pattern and begin to efficiently pull the data into the processor, thus reducing data access latency costs.


### Memory escape

In Go, variables with a fully known life cycle are allocated on the stack for efficiency.
If a variable's life cycle is not entirely predictable, it 'escapes' and is allocated on the heap for proper memory management.

Typical situations that can cause variables to escape onto the heap:

- Returning a local variable pointer within a method may extend its life cycle beyond the stack, causing a stack overflow

- Sending a pointer or a value with a pointer to a channel makes it challenging for the compiler to predict when
  the variable will be released, as the receiving goroutine is unknown at compile time

- Storing a pointer or value with a pointer on a slice, like []*string, leads to the slice's contents escaping to the heap,
  even if the array behind it is initially allocated on the stack.
  
- If the slice's capacity is exceeded during appending, reallocation occurs on the heap

- Calling a method on an interface type, such as invoking methods on ```io.Reader```, dynamically dispatches the method at runtime.
  This causes the value and the storage behind the slice to escape, resulting in heap allocation.

### Cache line

```
func RowTraverse() int {
    var ctr int
    for row := 0; row < rows; row++
        for col := 0; col < cols; col++ {
            if matrix[row][col] == 0xFF {
                ctr++
            }
        }
    }
    return ctr
}
```

Row traverse will have the best performance because it walks through memory,
cache line by connected cache line, which creates a predictable access pattern.

Cache lines can be prefetched and copied into the L1 or L2 cache before the data is needed.

```
func ColumnTraverse() int {
    var ctr int
    for col := 0; col < cols; col++ {
        for row := 0; row < rows; row++ {
            if matrix[row][col] == 0xFF {
                ctr++
            }
        }
    }
    return ctr
}
```

Column Traverse is the worst by an order of magnitude because this access pattern
crosses over OS page boundaries on each memory access. This causes no
predictability for cache line prefetching and becomes essentially random access
memory

```
func LinkedListTraverse() int {
    var ctr int
    d := list
    for d != nil {
        if d.v == 0xFF {
            ctr++
        }
        d = d.p
    }
    return ctr
}
```

The linked list is twice as slow as the row traversal mainly because there are cache
line misses but fewer TLB (Translation Lookaside Buffer) misses. A bulk of the
nodes connected in the list exist inside the same OS pages.

BenchmarkLinkListTraverse-16 128  28738407  ns/op
BenchmarkColumnTraverse-16   30   126878630 ns/op
BenchmarkRowTraverse-16      310  11060883  ns/op


# Translation Lookaside Buffer (TLB)

OS shares physical memory by breaking the physical memory into pages and mapping pages to virtual memory for any given running program. Each OS can decide the size of a page, but 4k, 8k, 16k are reasonable and common sizes.

The TLB is a small cache inside the processor that helps to reduce latency on translating a virtual address to a physical address within the scope of an OS page and offset inside the page.

TLB cache miss can cause large latencies because now the hardware has to wait for the OS to scan its page table to locate the right page for the virtual address


### Map keys

Slice is a good example of a type that can’t be used as a key. Only values that can
be run through the hash function are eligible. A good way to recognize types that
can be a key is if the type can be used in a comparison operation. I can’t compare
two slice values.



### Slices

```
aa := make([]string, 0)
aa = append(aa, "a")
aa = append(aa, "b")
aa = append(aa, "c")
aa = append(aa, "d") // len: 4 cap: 4
```

small capacity
in this case ```append``` creates a new backing array (doubling or growing by 25%)
and then copies the values from the `old` array into the `new` one

```
aa = append(aa, "e") // len: 5 cap: 8
```

Slices offer the capability to prevent additional copies and heap allocations of the underlying array

```
slice1 := []string{"A", "B", "C", "D", "E"} // len: 5 cap: 5
                              ^
                              |
                              pointer on that position

slice2 := slice1[2:4]                       // len: 2 cap: 3
```

```slice2``` only allows me to access the elements at index 2 and 3 (C and D) of the original slice’s backing array. The length of slice2 is 2 and not 5 like in slice1 and the capacity is 3 since there are now 3 elements from that pointer position.

