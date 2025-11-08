package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchRetriesUntilSuccess(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch n := atomic.AddInt32(&attempts, 1); {
		case n < 3:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(Item{ID: "id-1", Data: "payload"}); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := &http.Client{Timeout: time.Second}

	item, err := fetch(ctx, client, srv.URL)
	if err != nil {
		t.Fatalf("fetch returned error: %v", err)
	}

	if item.ID != "id-1" || item.Data != "payload" {
		t.Fatalf("unexpected item: %+v", item)
	}

	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts, got %d", got)
	}
}

func TestFetchSkipsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := &http.Client{Timeout: time.Second}

	_, err := fetch(ctx, client, srv.URL)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}

	if err.Error() != "skip404" {
		t.Fatalf("expected skip404 error, got %v", err)
	}
}

func TestIngestStoresEachItem(t *testing.T) {
	type response struct {
		status int
		item   Item
	}

	responses := map[string]response{
		"/one":   {status: http.StatusOK, item: Item{ID: "one", Data: "first"}},
		"/two":   {status: http.StatusOK, item: Item{ID: "two", Data: "second"}},
		"/three": {status: http.StatusOK, item: Item{ID: "three", Data: "third"}},
	}

	var mu sync.Mutex
	seen := make(map[string]int)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen[r.URL.Path]++
		mu.Unlock()

		resp, ok := responses[r.URL.Path]
		if !ok {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(resp.status)
		if resp.status >= 200 && resp.status < 300 {
			if err := json.NewEncoder(w).Encode(resp.item); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
	}))
	defer srv.Close()

	urls := []string{
		srv.URL + "/one",
		srv.URL + "/two",
		srv.URL + "/three",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	store := NewStore()
	if err := Ingest(ctx, urls, 2, store); err != nil {
		t.Fatalf("ingest returned error: %v", err)
	}

	if got := store.Len(); got != len(urls) {
		t.Fatalf("expected %d items in store, got %d", len(urls), got)
	}

	assert.Equal(t, len(urls), store.Len())

	for path, count := range seen {
		if count == 0 {
			t.Fatalf("expected path %s to be requested at least once", path)
		}
	}
}
