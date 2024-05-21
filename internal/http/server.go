package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zaidsasa/xbankapi/internal/logger"
)

const (
	httpServerReadHeaderTimeout = 1 * time.Second
	timeoutContextDuration      = 1 * time.Second
)

type (
	Handler interface {
		Register(mux *http.ServeMux)
	}

	Server struct {
		logger   logger.Logger
		handlers []Handler
	}
)

// NewServer returns a new Server.
func NewServer(logger logger.Logger, handlers ...Handler) *Server {
	return &Server{
		logger:   logger,
		handlers: handlers,
	}
}

// Start serving with the provided address.
func (s *Server) Start(ctx context.Context, addr string) error {
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

	errChan := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("failed to serve server: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}

	s.logger.Info("shutting down server, please wait...")

	ctx, cancelFunc := context.WithTimeout(ctx, timeoutContextDuration)
	defer cancelFunc()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	s.logger.Info("a graceful shutdown")

	return nil
}
