package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chutified/url-shortener/config"
	"github.com/chutified/url-shortener/controller"
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
// and server structure based on the given configuration + manage routings and endpoints.
func (s *server) Set(ctx context.Context, cfg *config.Config) error {

	// set timeout
	s.srvTimeOut, _ = time.ParseDuration(cfg.SrvTimeOut)

	// initialize handler
	s.h = controller.NewHandler()
	err := s.h.InitDataService(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("can not init handler's data service: %w", err)
	}

	// get handler with routings applied
	r := s.h.GetHTTPHandler()

	// create a server
	s.srv = &http.Server{
		Addr:              cfg.Addr(),
		Handler:           r,
		ReadTimeout:       500 * time.Millisecond,
		ReadHeaderTimeout: 300 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
	}

	return nil
}

// Run starts the server.
func (s *server) Run() error {

	// run server
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server can not be launched: %w", err)
	}

	return nil
}

// Stop stop the server.
func (s *server) Stop() error {

	// stop server
	ctx, cancel := context.WithTimeout(context.Background(), s.srvTimeOut)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("forced shutdown: %w", err)
	}

	return nil
}

// Close closes all open connections and services.
func (s *server) Close() error {

	// close handler
	err := s.h.CloseHandler()
	if err != nil {
		return fmt.Errorf("unsuccessful handler's closure: %w", err)
	}

	return nil
}
