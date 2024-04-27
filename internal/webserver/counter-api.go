package webserver

import (
	"fmt"
	"net/http"
)

const CounterHandlerPath = `/counter`

func (s *Server) CounterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.requestCounter.AddRequest()
		if err != nil {
			s.logger.Error(`failed to add a request`, `error`, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		count := s.requestCounter.CountRequests()

		w.Header().Set(`Content-Type`, `application/json`)
		_, err = w.Write([]byte(fmt.Sprintf("Total requests in the last 60 seconds: %d", count)))
		if err != nil {
			s.logger.Error(`failed to write response body`, `error`, err.Error())
		}
	})
}
