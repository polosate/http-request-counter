package webserver

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"simplesurance-test-task/internal/limiter"
)

func TestWithParallelLimiter(t *testing.T) {
	maxParallelRequest := 5
	parallelLimiter := limiter.NewParallelRateLimiter(maxParallelRequest)

	srv, err := New(
		``,
		``,
		nil,
		nil,
		parallelLimiter,
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	srv.httpMux.Handle("/testAPI", srv.withParallelLimiter(srv.testApi()))

	server := httptest.NewServer(srv.httpMux)
	defer server.Close()

	makeRequest := func() int {
		resp, err := http.Get(server.URL + "/testAPI")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()
		return resp.StatusCode
	}

	var wg sync.WaitGroup
	results := make(chan int, maxParallelRequest+1)

	for i := 0; i < maxParallelRequest+1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- makeRequest()
		}()
	}

	wg.Wait()
	close(results)

	allowed := 0
	tooManyRequests := 0

	for result := range results {
		if result == http.StatusOK {
			allowed++
		} else if result == http.StatusTooManyRequests {
			tooManyRequests++
		} else {
			t.Errorf("Unexpected status code: %d", result)
		}
	}

	if allowed != 5 {
		t.Errorf("Expected 2 allowed requests, got %d", allowed)
	}
	if tooManyRequests != 1 {
		t.Errorf("Expected 1 too many requests, got %d", tooManyRequests)
	}

	// Wait for the first request to release the slot
	time.Sleep(3 * time.Second)

	resp, err := http.Get(server.URL + "/testAPI")
	if err != nil {
		t.Fatalf("Failed to make third request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}
}

func (s *Server) testApi() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	})
}
