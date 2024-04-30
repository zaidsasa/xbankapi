package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const httpServerReadHeaderTimeout = 3 * time.Second

type (
	Handler interface {
		Register(mux *http.ServeMux)
	}

	Server struct {
		logger   *slog.Logger
		handlers []Handler
	}
)

// NewServer returns a new Server.
func NewServer(logger *slog.Logger, handlers ...Handler) *Server {
	return &Server{
		logger:   logger,
		handlers: handlers,
	}
}

// Start serving with the provided address.
func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	for _, handler := range s.handlers {
		handler.Register(mux)
	}

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: httpServerReadHeaderTimeout,
		Handler:           mux,
	}

	s.logger.Info("server started", "address", addr)

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to serve server :%w", err)
	}

	return nil
}
