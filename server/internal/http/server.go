package http

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string, router http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
