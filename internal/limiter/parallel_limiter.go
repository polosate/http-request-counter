package limiter

import "sync"

type ParallelRateLimiter struct {
	Semaphore chan struct{}
}

func NewParallelRateLimiter(maxParallelRequests int) *ParallelRateLimiter {
	return &ParallelRateLimiter{
		Semaphore: make(chan struct{}, maxParallelRequests),
	}
}

func (l *ParallelRateLimiter) Allow() bool {
	select {
	case l.Semaphore <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *ParallelRateLimiter) Release() {
	<-l.Semaphore
}

type LimiterManager struct {
	limiters map[string]*ParallelRateLimiter
	mu       sync.RWMutex
}

func NewLimiterManager() *LimiterManager {
	return &LimiterManager{
		limiters: make(map[string]*ParallelRateLimiter),
	}
}

func (m *LimiterManager) AddLimiter(path string, limiter *ParallelRateLimiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[path] = limiter
}

func (m *LimiterManager) GetLimiter(path string) *ParallelRateLimiter {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.limiters[path]
}
