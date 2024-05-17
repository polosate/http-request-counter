package limiter

type ParallelRateLimiter struct {
	Semaphore chan struct{}
}

func NewParallelRateLimiter(maxParallelRequests int) *ParallelRateLimiter {
	return &ParallelRateLimiter{
		Semaphore: make(chan struct{}, maxParallelRequests),
	}
}
