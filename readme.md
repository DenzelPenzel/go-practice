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
  ```
  package sync

  type Mutex struct {
    state int32
    sema  uint32
  }
  ```
  - sync.Mutex
  - sync.RWMutex - multiple readers acquire the lock simultaneously but only one writer at a time
  
  - Mutexes using low-level atomic operations and OS-level synchronization primitives 
    (like futexes on Linux) provided by the Go runtime
  
  - If the mutex is already <b>locked</b>, the runtime adds the goroutine to a <b>waiting queue</b> and blocks it
  
  - Atomic operations are supported through specific CPU instructions like 
    ```test-and-set```, ```compare-and-swap```, or ```fetch-and-add```
  
  - OS maintains a queue of threads that are waiting for the mutex
  
  - When mutex is released, the OS wakes up one or more threads from the queue and allows them to try to acquire the mutex
  

- WaitGroup
  ```
  type WaitGroup struct {
	  noCopy noCopy

	  state atomic.Uint64
	  sema  uint32
  }

  type noCopy struct{}

  func (*noCopy) Lock()   {}
  func (*noCopy) Unlock() {}
  ```

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

- Channels are <b>concurrency primitive</b> and can be used to synchronize operations between G

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

### If you create a channel, send data into it, but never read from it

- unbuffered channel: program will block indefinitely
- buffered channel: if buffer become full and no receiver, program will block


### Allocation

Allocation on the ```stack``` typically happens for local variables

Allocation on the ```heap``` happens when the compiler determines 
the variable escapes the local scope (e.g., returned from a function, stored in a heap-allocated object).

```
func createInt() *int {
  x := 42
  return &x // x escapes to the heap
}
```

### Allocating Structs with ```new```:

Using the ```new``` allocates memory on the heap.

```
  type MyStruct struct{ A int }
  p := new(MyStruct) // Heap allocation
```

Notes:
  - Avoid unnecessary pointers which may lead to heap allocations.
   ```
    func processInt(x int) {
      // x passed by value, no heap allocation
    }
  ``` 

  - Slicing a large array can keep the entire array in memory. Create a smaller copy if needed.
    ```
      largeArray := [1000]int{}
      smallSlice := make([]int, 10)
      copy(smallSlice, largeArray[:10]) // Avoid holding onto largeArray
    ```

  - Use ```sync.Pool``` for Reusable Objects(avoid frequent allocations and garbage collection)
    ```
    var bufferPool = sync.Pool{
      New: func() interface{} {
          return make([]byte, 1024)
      },
    }
    func process() {
      buf := bufferPool.Get().([]byte)
      // Use buf
      bufferPool.Put(buf) // Reuse buffer
    }

  - Use Go’s profiling tools ```pprof``` ```trace```

    ```


If the function can be inlined, the ownership of the construction moves up to the calling function

```
func fn() {
  // Before inlining
  input := bytes.NewReader(data) // <- Original Call
   
  // After inlining
  input := &bytes.Reader{ buf: data } // <- After Inlining Optimization
}
```

### GC mark assist

GC activates the GC ```mark assist``` to speed up the marking process.

GC ```mark assist``` in Go helps dynamically adjusting the resources allocated to the garbage collector, during the mark phase to keep up with the pace of object allocation

Could be the situation when all the Goroutines in ```mark assist``` to help slowdown allocations and get the initial GC finished.

### Mark assist algo

Core ideas:
  - Mark Phase: identifies all the live (reachable) objects in the heap
  
  - Sweep Phase: destroy memory occupied by objects that were not marked as live
  
  - GC recursively to find all reachable objects. 
  
  - Go use "tri-color marking"
    - White: Objects that haven’t been visited yet (initially all objects) -> Remove this objects later
    - Gray: Objects that have been marked but whose references haven’t been scanned
    - Black: Objects that have been fully processed (marked and their references scanned)

### Performance notes

- If you notice a lot of time spent in the ```runtime.mallocgc``` function, it suggests that the program may be making too many small memory allocations

- If you're spending a significant amount of time managing channel operations, ```sync.Mutex``` code, or other synchronization elements in your program, it's likely facing contention issues. 

To improve performance, think about restructuring the program to reduce the frequent access of shared resources. 

Common techniques for this include techniques:
- ```sharding/partitioning``` 
- ```buffering/batching``` 
- ```copy-on-write``` 

- If your program spends a significant amount of time in ```syscall.Read/Write```, it might be doing too many small reads and writes. 
  Using ```bufio``` wrappers around ```os.File``` or ```net.Conn``` can be helpful in this situation

- If your program is spending a lot of time in the ```GC (Garbage Collection)``` component, it could be because it's either creating too many temporary objects or because the heap size is too small, leading to frequent garbage collections:
  
    - Large objects impact memory usage and GC pacing, whereas numerous small allocations affect marking speed

    - Combine values into larger ones to reduce memory allocations and alleviate pressure on the garbage collector, resulting in faster garbage collections

    - Values without pointers aren't scanned by the garbage collector. Eliminating pointers from actively used values can enhance garbage collection efficiency

### Inlining optimization

Inlining is an optimization technique used by compilers, including the Go compiler, to improve the performance of a program.

Basic idea behind inlining is to replace a function call with the actual body of the function

This can eliminate the overhead associated with the function call, such as stack manipulation and jump instructions, thereby making the program faster.


Cons:
  - Increased Binary Size
  - Larger binaries can negatively impact CPU cache performance

Viewing Inlining Decisions
  - use the ```-gcflags``` compiler flag with ```go build``` or ```go test```


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
- Minimize  context switches  
- Keep local queue per P, and global queue
- If local queues of P is full, the scheduler uses a global run queue
- If the load is imbalanced, the scheduler may migrate goroutines between P to better distribute work

### Scheduling Decisions
- Preemption - ensure that no single goroutine can monopolize a P

### Blocking Operations
When a goroutine performs a blocking operation (like I/O or a system call), the ```M``` executing that ```G``` may block. The scheduler can detach the ```P``` from the blocked ```M``` and attach it to another ```M``` to continue running other goroutines, ensuring that the system remains responsive. Once the blocking operation completes, the original ```M``` can be reused

<b>Idle M Handling</b>: 
  If an ```M``` becomes idle and there are more ```Ms``` than ```Ps```, the Go runtime may park the idle ```M``` (put it to sleep), reducing resource usage


### Parallelism
- GOMAXPROCS set to the number of cores of the current machine
- 4 active operating system threads will be created on a 4 core machine
- If GOMAXPROCS is set to more than 1, the Go runtime can execute multiple goroutines in parallel 
  on multiple processors. This allows the Go runtime to utilize multiple CPU cores


### CPU caches

Prefetcher attempts to predict what data is needed 
before instructions request the data so it’s already present in either the L1 or L2 cache.

Program can read/write a byte of memory as the smallest unit of memory access.

Caching systems granularity is 64 bytes. 
This 64 byte block of memory is called a cache line.

Prefetcher works best when the instructions being executed create predictable access patterns to memory. 
One way to create a predictable access pattern to memory is to construct a contiguous block of memory and then iterate over that memory performing a linear traversal with a predictable stride.

Array is data structure with predictable access patterns. 
However, the slice is the most impoдrtant data structure in Go. Slices in Go use an array underneath.

Prefetcher will pick up predictable data access pattern and begin to efficiently pull the data into the processor, thus reducing data access latency costs.


### Memory escape

In Go, variables with a fully known life cycle are allocated on the stack for efficiency.
If a variable's life cycle is not entirely predictable, it 'escapes' and is allocated on the heap for proper memory management.

Typical situations that can cause variables to escape onto the heap:

- Returning a local variable pointer within a method may extend its life cycle beyond the stack, causing a stack overflow

- Sending a pointer or a value with a pointer to a channel makes it challenging for the compiler to predict when
  the variable will be released, as the receiving goroutine is unknown at compile time

- Storing a pointer or value with a pointer on a slice, like []*string, leads to the slice's contents escaping to the heap, even if the array behind it is initially allocated on the stack.
  
- If the slice's capacity is exceeded during appending, reallocation occurs on the heap

- Calling a method on an interface type, such as invoking methods on ```io.Reader```, dynamically dispatches the method at runtime. This causes the value and the storage behind the slice to escape, resulting in heap allocation.

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

OS shares physical memory by breaking the physical memory into pages 
and mapping pages to virtual memory for any given running program. 

Each OS can decide the size of a page, but 4k, 8k, 16k are reasonable and common sizes.

The TLB is a small cache inside the processor that helps to reduce latency on translating a virtual address to a physical address within the scope of an OS page and offset inside the page.

TLB cache miss can cause large latencies because now the hardware has to wait for the OS to scan its page table to locate the right page for the virtual address


### Map keys

Slice is a good example of a type that can’t be used as a key. Only values that can
be run through the hash function are eligible. A good way to recognize types that
can be a key is if the type can be used in a comparison operation. I can’t compare
two slice values.


### Slices vs Array


### Slices

- Slice descriptor contains a pointer to an underlying array, along with length and capacity

- Pointer to a slice (&aaa) -> ref to the slice descriptor, not the underlying array
  ```
    aaa := []int{7, 8, 9}
    for i, v := range *(&aaa) {
      fmt.Println(i, v)
    }
  ```

### Small capacity of slices

```
aa := make([]string, 0)
aa = append(aa, "a")
aa = append(aa, "b")
aa = append(aa, "c")
aa = append(aa, "d") // len: 4 cap: 4

// small capacity
// in this case ```append``` creates a new backing array (doubling or growing by 25%)
// and then copies the values from the `old` array into the `new` one
aa = append(aa, "e") // len: 5 cap: 8
```

In this case ```append``` creates a new backing array (doubling or growing by 25%)
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


### Race condition 
- Two or more goroutines access shared data concurrently
- Issues:
  - Unpredictable behavior
  - Data corruption
  - Hard to Debug
  
  

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


### Run Bench

```
// generate pprof and binary
# go test -bench . -benchmem -memprofile p.out -gcflags -m=2
# go test -bench . -benchtime 3s -benchmem -memprofile p.out
# go test -bench . -benchtime 3s -benchmem -cpuprofile p.out -gcflags -m=2

// run tools
# go tool pprof -noinlines p.out
# go tool pprof --noinlines  memcpu.test p.out

// Profiling commands
# go tool pprof -http :8080 stream.test p.out
# press ```o```
# list <func_name>
# weblist <func_name>
```


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