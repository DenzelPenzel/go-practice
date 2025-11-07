package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Item struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Item
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Item),
	}
}

func (s *Store) Put(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[item.ID] = item
}

func (s *Store) Get(id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.data[id]
	return item, ok
}

func (s *Store) Len() int {
	s.mu.RLock()
	n := len(s.data)
	s.mu.RUnlock()
	return n
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
	return nil
}

func retry(ctx context.Context, attempts int, base time.Duration, fn func() error) error {
	d := base
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i == attempts-1 {
			break
		}
		j := time.Duration(rand.Int63n(int64(d / 2)))
		select {
		case <-time.After(d + j):
		case <-ctx.Done():
			return ctx.Err()
		}
		d *= 2
	}
	return err
}

func fetchOnce(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func fetch(ctx context.Context, client *http.Client, url string) (Item, error) {
	var resp *http.Response

	err := retry(ctx, 3, 150*time.Millisecond, func() error {
		r, err := fetchOnce(ctx, client, url)
		if err != nil {
			return err
		}
		resp = r
		if resp.StatusCode == http.StatusNotFound {
			// drains HTTP response body, close body without reading it
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			return errors.New("skip404")
		}
		if resp.StatusCode >= 500 {
			// drains HTTP response body, close body without reading it
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			return fmt.Errorf("retryable %d", resp.StatusCode)
		}
		return nil
	})

	if err != nil {
		if err.Error() == "skip404" {
			return Item{}, err
		}
		return Item{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// non-404, non-2xx -> fatal
		return Item{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var it Item
	if err := json.NewDecoder(resp.Body).Decode(&it); err != nil {
		return Item{}, err
	}

	return it, nil
}

func Ingest(ctx context.Context, urls []string, concurrency int, store *Store) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	g, ctx := errgroup.WithContext(ctx)

	// 1.12+ gives each iteration of a range loop its own copy of the element
	for _, u := range urls {
		g.Go(func() error {
			it, err := fetch(ctx, client, u)
			if err != nil {
				// tolerate 404 marker
				if err.Error() == "skip404" {
					return nil
				}
				return err
			}
			store.Put(it)
			return nil
		})
	}
	return g.Wait()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store := NewStore()
	urls := []string{
		"https://httpbin.org/json",          // not Item shape; will error (demo)
		"https://httpbin.org/status/404",    // tolerated
		"https://httpbin.org/anything?ok=1", // returns JSON; still may not match Item
	}

	err := Ingest(ctx, urls, 10, store)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}
