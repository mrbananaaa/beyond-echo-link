package http

import "net/http"

type Server struct {
	server  *http.Server
	handler http.Handler
}

func New(
	addr string,
	handler http.Handler,
) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
