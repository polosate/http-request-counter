package webserver

import (
	"net/http"
)

func (s *Server) withLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(
			`Got request`,
			`url`, r.URL.String(),
			`header`, r.Header,
		)
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func (s *Server) withParallelLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := s.limiterManager.GetLimiter(r.URL.Path)
		if limiter == nil {
			http.Error(w, "Rate limiter not configured for this path", http.StatusInternalServerError)
			return
		}

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		defer limiter.Release()

		next.ServeHTTP(w, r)
	})
}
