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

### Goroutine

- Each Goroutine is given its own block of memory called a stack
- Each stack starts out as a 2048 byte (2k) allocation
- Func is called -> allocation of the stack space to execute func -> this block of memory is called a frame
- Size of a frame for a given function is calculated at compile time
- If the compiler doesn’t know the size of a value at compile time, the
  value has to be constructed on the heap

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

### Map keys

Slice is a good example of a type that can’t be used as a key. Only values that can
be run through the hash function are eligible. A good way to recognize types that
can be a key is if the type can be used in a comparison operation. I can’t compare
two slice values.



