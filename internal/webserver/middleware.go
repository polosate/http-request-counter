package webserver

import "net/http"

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
