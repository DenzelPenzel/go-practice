/*

Situation: In a web server with high levels of concurrency,
there is a need to control frequent access by IP addresses.

The scenario involves simulating 100 different IPs accessing the server simultaneously,
with each IP making 1000 requests. The condition is that each IP should be allowed only
one access within a three-minute window.

Your task is to adapt the provided code to achieve this and ensure
that it successfully outputs "Success: 100"

*/

package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Ban struct {
	visitIPs map[string]time.Time
	lock     sync.Mutex
}

func NewBan(ctx context.Context) *Ban {
	b := &Ban{
		visitIPs: make(map[string]time.Time),
	}

	go func() {
		timer := time.NewTimer(time.Minute * 1)

		for {
			select {
			case <-timer.C:
				b.lock.Lock()
				for k, v := range b.visitIPs {
					if time.Now().Sub(v) >= time.Minute*1 {
						delete(b.visitIPs, k)
					}
				}
				b.lock.Unlock()
				timer.Reset(time.Minute * 1)

			case <-ctx.Done():
				return
			}
		}
	}()

	return b
}

func (b *Ban) visit(ip string) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.visitIPs[ip]; ok {
		return true
	}
	b.visitIPs[ip] = time.Now()
	return false
}

func main() {
	// When values are altered concurrently on multiple CPU cores,
	// there's a possibility that the integer value may become unsynchronized in exceptional situations.
	// To address this, it's crucial to ensure that modifications to the integer value are performed atomicall
	success := int64(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ban := NewBan(ctx)

	wait := &sync.WaitGroup{}
	wait.Add(1000 * 100)

	for i := 0; i < 1000; i++ {
		for j := 0; j < 100; j++ {
			go func(j int) {
				defer wait.Done()
				ip := fmt.Sprintf("192.168.1.%d", j)
				if !ban.visit(ip) {
					atomic.AddInt64(&success, 1)
				}
			}(j)
		}
	}

	wait.Wait()

	fmt.Println("success:", success)
}
