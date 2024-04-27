package webserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"simplesurance-test-task/internal/counter"
)

type Server struct {
	httpLis net.Listener
	httpMux *http.ServeMux
	httpSrv *http.Server

	logger *slog.Logger

	requestCounter *counter.RequestCounter
}

func New(
	addr string,
	port string,
	logger *slog.Logger,
	requestCounter *counter.RequestCounter,
) (*Server, error) {
	if addr == `` {
		addr = `:` + port
	} else {
		addr = addr + `:` + port
	}

	lis, err := net.Listen(`tcp`, addr)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		httpLis: lis,
		httpMux: http.NewServeMux(),
		httpSrv: &http.Server{
			Addr:         addr,
			WriteTimeout: responseTimeout,
		},
		logger:         logger,
		requestCounter: requestCounter,
	}

	srv.httpMux.Handle(
		CounterHandlerPath,
		srv.withLogRequest(srv.CounterHandler()),
	)

	srv.httpSrv.Handler = srv.httpMux

	return srv, nil
}

func (s *Server) Serve() error {
	return s.httpSrv.Serve(s.httpLis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

const (
	responseTimeout = 1 * time.Minute
)
