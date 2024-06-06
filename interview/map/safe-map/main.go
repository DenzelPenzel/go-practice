/*
Implement blocking read and concurrency-safe map

Questions:
	- How does MAP in GO realize that the key does not exist?

The get operation waits until the key exists or times out to ensure concurrency safety,
and the following interface needs to be implemented:
*/

package main

import (
	"log"
	"sync"
	"time"
)

type Mp interface {
	// Storing key/value pairs will activate a wake-up, if the coroutine attempting to read the key is in a suspended state
	// This method operates without blocking and is capable of immediate execution and return
	Set(key string, val interface{})

	// Reading a key, if the key does not exist, it blocks
	// Waiting for the key to become available or for a timeout to occur
	Get(key string)
}

type Map struct {
	mapping map[string]*record
	rmx     *sync.RWMutex
}

type record struct {
	ch      chan struct{}
	value   interface{}
	isExist bool
}

func (m *Map) Set(key string, val interface{}) {
	m.rmx.Lock()
	defer m.rmx.Unlock()

	item, ok := m.mapping[key]

	if !ok {
		m.mapping[key] = &record{
			value:   val,
			isExist: true,
		}
		return
	}

	item.value = val

	if !item.isExist {
		if item.ch != nil {
			close(item.ch)
			item.ch = nil
		}
	}
}

func (m *Map) Get(key string, timeout time.Duration) interface{} {
	m.rmx.RLock()
	r, ok := m.mapping[key]

	if ok {
		if r.isExist {
			m.rmx.RUnlock()
			return r.value
		} else {
			m.rmx.RUnlock()
			log.Println("waiting... key ->", key)
			select {
			case <-r.ch:
				return r.value
			case <-time.After(timeout):
				log.Println("timeout -> ", key)
				return nil
			}
		}
	} else {
		m.rmx.RUnlock()
		// exclusive access to data
		m.rmx.Lock()
		r = &record{
			ch:      make(chan struct{}),
			isExist: false,
		}

		m.mapping[key] = r
		m.rmx.Unlock()

		log.Println("waiting... key ->", key)

		select {
		case <-r.ch:
			return r.value
		case <-time.After(timeout):
			log.Println("timeout -> ", key)
			return nil
		}
	}
}

func main() {
	mp := Map{
		mapping: make(map[string]*record),
		rmx:     &sync.RWMutex{},
	}

	for i := 0; i < 10; i++ {
		go func() {
			val := mp.Get("key", time.Second*5) // blocking, while timeout or
			log.Println("done... get ->", val)
		}()
	}

	time.Sleep(time.Second * 3)

	for i := 0; i < 10; i++ {
		go func(val int) {
			mp.Set("key", val)
		}(i)
	}

	time.Sleep(time.Second * 30)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
