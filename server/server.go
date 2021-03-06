package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/chutommy/url-shortener/config"
	"github.com/chutommy/url-shortener/controller"
)

const (
	readTimeOut     = 500
	readHearTimeout = 300
	writeTimeout    = 500
)

// Server manages the whole web service runtime.
type Server interface {
	Set(context.Context, *config.Config) error
	Run() error
	Stop() error
	Close() error
}

// server implements Server interface.
type server struct {
	h          controller.Handler
	srv        *http.Server
	srvTimeOut time.Duration
}

// NewServer is a constructor of the server.
func NewServer() Server {
	return &server{}
}

// Set prepares server to run. Set creates under the hood a new database connection
// and server structure based on the given configuration + manage routing and endpoints.
func (s *server) Set(ctx context.Context, cfg *config.Config) error { // set timeout
	s.srvTimeOut, _ = time.ParseDuration(cfg.SrvTimeOut)

	// set handler
	if err := s.setHandler(ctx, cfg); err != nil {
		return fmt.Errorf("failed to set handler: %w", err)
	}

	// set server
	s.setServer(cfg)

	return nil
}

// Run starts the server.
func (s *server) Run() error {
	// run server
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server can not be launched: %w", err)
	}

	return nil
}

// Stop stops the server.
func (s *server) Stop() error {
	// stop server
	ctx, cancel := context.WithTimeout(context.Background(), s.srvTimeOut)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("a forced shutdown failed: %w", err)
	}

	return nil
}

// Close closes all open connections and services.
func (s *server) Close() error { // close handler
	err := s.h.CloseHandler()
	if err != nil {
		return fmt.Errorf("an unsuccessful handler's closure: %w", err)
	}

	return nil
}

// setHandler initializes handler's data service and sets it for the server.
func (s *server) setHandler(ctx context.Context, cfg *config.Config) error {
	// initialize handler
	s.h = controller.NewHandler()

	err := s.h.InitDataService(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("can not init handler's data service: %w", err)
	}

	return nil
}

// setServer constructs a server.
func (s *server) setServer(cfg *config.Config) {
	// get handler with routing applied
	r := s.h.GetHTTPHandler()

	// create a server
	s.srv = &http.Server{
		Addr:              cfg.Addr(),
		Handler:           r,
		ReadTimeout:       readTimeOut * time.Millisecond,
		ReadHeaderTimeout: readHearTimeout * time.Millisecond,
		WriteTimeout:      writeTimeout * time.Millisecond,
	}
}
