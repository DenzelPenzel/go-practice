package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	test1()
}

// 1 Mutex
// error due to a deadlock
func test1() {
	var mu sync.Mutex
	var chain string

	// deadlock, the lock is held by the A(), function C() function is trying to acquire the same lock without releasing it
	C := func() {
		mu.Lock() // <---- deadlock
		defer mu.Unlock()
		chain = chain + " --> C"
	}

	B := func() {
		// mu.Lock() // ---> fix issue with deadlock

		// update without release the lock
		chain = chain + " --> B"

		// mu.Unlock() // ---> fix issue with deadlock
		C()
	}

	A := func() {
		mu.Lock()
		defer mu.Unlock()
		chain = chain + " --> A"

		//mu.Unlock() // ---> fix issue with deadlock

		B()
	}

	chain = "main"
	A()
	fmt.Println(chain)
}

// ==========================================================================================================
// 2 RWMutex
// deadlock issue
func test2() {
	var mu sync.RWMutex
	var count int

	// trying to acquire a read lock in function C() while holding a read lock in function A()
	C := func() {
		mu.RLock()
		defer mu.RUnlock()
	}

	B := func() {
		time.Sleep(5 * time.Second)
		// mu.RUnlock() // ---> to fix issue need to release the read lock before calling C()
		C()
	}

	A := func() {
		mu.RLock()
		defer mu.RUnlock()
		B()
	}

	go A()
	time.Sleep(time.Second * 2)
	mu.Lock()
	defer mu.Unlock()
	count++
	fmt.Println(count)
}

// ==========================================================================================================
// 3. WaitGroup
// race condition
func test3() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		time.Sleep(time.Millisecond)
		wg.Done()
		wg.Add(1) // ---> race condition
	}()

	wg.Wait()
	// wg.Add(1) correct version
}

// ==========================================================================================================
// 4 Double check to implement singleton
// In multi-core CPUs, the CPU cache will cause variable values ​​in multiple cores to be out of sync

// atomic vs Mutex
// atomic - better performance because can directly manipulate memory without introducing locks
// atomic better for the simple operations

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	o.m.Lock()
	defer o.m.Unlock()

	/*
		correct version and not need to use Lock here
		if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
			f()
		}
	*/

	if o.done == 0 {
		o.done = 1
		f()
	}
}

// ==========================================================================================================
// 5 Mutex
// deadlock issue
type MyMutex struct {
	count int
	sync.Mutex
}

func test5() {
	var mu MyMutex
	mu.Lock()
	// copying the variable after locking will also copy the lock status
	// so it mu actually already locked
	// if you lock it again, it will cause a deadlock
	var mu2 = mu
	mu.count++
	mu.Unlock()

	mu2.Lock() // lock status copied ---> deadlock
	mu2.count++
	mu2.Unlock()
	fmt.Println(mu.count, mu2.count)
}

// ==========================================================================================================
// 7 channel
func test7() {
	var ch chan int

	go func() {
		ch = make(chan int, 1)
		ch <- 1
		fmt.Println("write")
	}()
	go func() {
		time.Sleep(time.Second * 1)
		val := <-ch
		fmt.Printf("read %d\n", val)
	}()

	c := time.Tick(1 * time.Second)
	for range c {
		fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
	}
}

// ==========================================================================================================
// 8 channel

func test8() {
	ch := make(chan int, 1)
	var count int
	go func() {
		ch <- 1
	}()
	go func() {
		count++
		close(ch)
	}()
	<-ch
	fmt.Println(count)
}

// ==========================================================================================================
// 9 Map
func test9() {
	var m sync.Map
	m.LoadOrStore("a", 1)
	m.Delete("a")
	// fmt.Println(m.Len()) // m.Len undefined
}

// ==========================================================================================================
// 9 Map
func test10() {
	var c = make(chan int)
	var a int

	go func() {
		a = 1
		// goroutine is trying to read from the channel, but there's no other goroutine sending a value
		<-c // deadlock
	}()
	// value 0 is sent on channel goroutine has a chance to execute and read from the channel
	c <- 0
	fmt.Println(a)
}

func test10Final() {
	var c = make(chan int)
	var a int
	var wg sync.WaitGroup

	wg.Add(1)

	go fn(&wg, c, &a)

	c <- 0
	wg.Wait()
	fmt.Println(a)
}

func fn(wg *sync.WaitGroup, c chan int, a *int) {
	defer wg.Done()
	*a = 1
	<-c
}
