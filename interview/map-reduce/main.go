package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type MapFn[I any, O any] func(ctx context.Context, input I) (O, error)
type ReduceFn[O any, R any] func(R, O) R

func MapReduce[I any, O any, R any](
	ctx context.Context,
	inputs []I,
	mapFn MapFn[I, O],
	reduceInit R,
	reduceFn ReduceFn[O, R],
) (R, error) {
	workers := runtime.NumCPU()
	inCh := make(chan I)
	outCh := make(chan O, 128)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range inCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				o, err := mapFn(ctx, item)
				if err != nil {
					return
				}

				select {
				case outCh <- o:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(inCh)
		for _, it := range inputs {
			select {
			case inCh <- it:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(outCh)
	}()

	acc := reduceInit
	for {
		select {
		case <-ctx.Done():
			return acc, ctx.Err()
		case o, ok := <-outCh:
			if !ok {
				return acc, nil
			}
			acc = reduceFn(acc, o)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inputs := []int{1, 2, 3, 4, 5}
	mapFn := func(ctx context.Context, x int) (int, error) {
		time.Sleep(50 * time.Millisecond)
		return x * x, nil
	}
	reduceFn := func(acc, o int) int { return acc + o }

	sum, err := MapReduce(ctx, inputs, mapFn, 0, reduceFn)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("sum of squares:", sum)
}
