package main

import (
	"log"
	"sync"
	"time"
)

// Find the error with the following code

type Users struct {
	mapping map[string]int
	mx      *sync.Mutex
}

/*
	Although sync.Mutex is used for write locks, concurrent reading and writing of map is unsafe
	Map is a reference type
	When reading and writing concurrently, multiple coroutines access the same address through pointers, that is, access shared variables
	And at this time, there is competition between reading and writing resources at the same time
	The error message will be reported: "fatal error: concurrent map read and map write".


*/

func (u *Users) Add(name string, age int) {
	u.mx.Lock()
	defer u.mx.Unlock()
	u.mapping[name] = age
}

// locks are also needed in because it is only for reading. It is recommended to use read-write locks sync.RWMutex
func (u *Users) Get(name string) int {
	if age, ok := u.mapping[name]; ok {
		return age
	}
	return -1
}

func main() {
	user := Users{
		mapping: make(map[string]int),
		mx:      &sync.Mutex{},
	}

	for i := 0; i < 10; i++ {
		go func() {
			val := user.Get("key")
			log.Println("done... get ->", val)
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			user.Add("key", i)
		}(i)
	}

	time.Sleep(time.Second * 10)

}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
