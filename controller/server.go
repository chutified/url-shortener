package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chutified/url-shortener/config"
)

// Server manages the whole web service runtime.
type Server interface {
	Set(*config.Config) error
	Run() error
	Stop() error
}

// server implements Server interface.
type server struct {
	h   *handler
	srv *http.Server
}

// NewServer is a constructor of the server.
func NewServer() Server {
	return &server{}
}

// Set prepares server to run. Set creates under the hood a new database connection
// and server structure based on the given configuration + manage routings and endpoints.
func (s *server) Set(cfg *config.Config) error {

	// initialize handler
	s.h = newHandler()
	err := s.h.initDataService(cfg.DB)
	if err != nil {
		return fmt.Errorf("can not init handler's data service: %w", err)
	}

	// get handler with routings applied
	r := s.h.getHTTPHandler()

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
		return fmt.Errorf("unexpected server error: %w")
	}

	return nil
}

// Stop stop the server and close all open connections and services.
func (s *server) Stop() error {
	//TODO
}
