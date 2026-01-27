// Package httpserver implements HTTP server using net/http.
package httpserver

import (
	"ReviewTask/pkg/logger"
	"context"
	"errors"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// Server represents an HTTP server wrapper.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	HTTPServer *http.Server // Public field for handler/customization
	notify     chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New creates a new HTTP server instance.
// User must set HTTPServer.Handler before calling Start().
func New(l logger.Interface, opts ...Option) *Server {
	// Build configuration using functional options pattern
	cfg := defaultOptions()
	for _, opt := range opts {
		opt(cfg)
	}

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Ensure single server goroutine

	return &Server{
		ctx: ctx,
		eg:  group,
		HTTPServer: &http.Server{
			Addr:         cfg.address,
			ReadTimeout:  cfg.readTimeout,
			WriteTimeout: cfg.writeTimeout,
			// Handler intentionally nil - user must configure before Start()
		},
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.shutdownTimeout,
		logger:          l,
	}
}

// Start launches HTTP server in managed goroutine.
// Requires HTTPServer.Handler to be set prior to invocation.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.HTTPServer.ListenAndServe()

		// Propagate unexpected errors (exclude graceful shutdown signal)
		if err != nil && err != http.ErrServerClosed {
			s.notify <- err
		}

		close(s.notify) // Signal termination to listeners
		return nil      // Treat ErrServerClosed as successful termination
	})

	s.logger.Info("restapi server - Server - Started", "address", s.HTTPServer.Addr)
}

// Notify returns channel receiving startup/runtime errors.
// Channel closes when server terminates (gracefully or with error).
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown performs graceful server termination with timeout enforcement.
func (s *Server) Shutdown() error {
	var errs []error

	// Initiate HTTP server shutdown with timeout context
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.HTTPServer.Shutdown(shutdownCtx); err != nil {
		if !errors.Is(err, context.Canceled) {
			s.logger.Error(err, "restapi server - Server - Shutdown - HTTPServer.Shutdown")
			errs = append(errs, err)
		}
	}

	// Wait for server goroutine completion
	if err := s.eg.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			s.logger.Error(err, "restapi server - Server - Shutdown - errgroup.Wait")
			errs = append(errs, err)
		}
	}

	s.logger.Info("restapi server - Server - Shutdown completed")
	return errors.Join(errs...)
}
