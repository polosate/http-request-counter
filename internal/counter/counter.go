package counter

import (
	"sync"
	"time"

	"simplesurance-test-task/internal/storage"
)

type RequestCounter struct {
	Requests map[time.Time]int
	mu       sync.Mutex
	storage  storage.Storage
}

func New(stor storage.Storage) (*RequestCounter, error) {
	rc := &RequestCounter{
		storage:  stor,
		Requests: make(map[time.Time]int),
	}

	// Load previous data if available
	data, err := rc.storage.Load()
	if err != nil {
		return rc, err
	}
	rc.Requests = data
	return rc, nil
}

func (rc *RequestCounter) AddRequest() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	now := time.Now().UTC().Truncate(time.Second)
	rc.Requests[now]++

	err := rc.storage.Save(rc.Requests)
	if err != nil {
		rc.Requests[now]--
		return err
	}
	return nil
}

func (rc *RequestCounter) CountRequests() int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	expirationTime := time.Now().UTC().Add(-60 * time.Second)
	count := 0
	for t, cnt := range rc.Requests {
		if t.Before(expirationTime) {
			delete(rc.Requests, t)
		} else {
			count += cnt
		}
	}
	return count
}
