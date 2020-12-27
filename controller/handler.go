package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chutified/url-shortener/config"
	"github.com/chutified/url-shortener/data"
)

// Handler is a handler interface of the controller
type Handler interface {
	GetHTTPHandler() http.Handler
	CloseHandler() error
	InitDataService(context.Context, *config.DB) error
}

// handler is the controller of the data service actions.
type handler struct {
	ds data.Service
}

// NewHandler returns an empty handler.
func NewHandler() Handler {
	return &handler{}
}

// InitDataService initializes handler's data service.
func (h *handler) InitDataService(ctx context.Context, dbCfg *config.DB) error {

	// create new data service
	h.ds = data.NewService()

	// initialize data service
	err := h.ds.InitDB(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize data service: %w", err)
	}

	return nil
}

// CloseHandler stops all active connections. closeHandler closes the data service.
// This function should not be called often (meant to be used only when the server
// is shutting down).
func (h *handler) CloseHandler() error {

	// close data service connection
	err := h.ds.StopDB()
	if err != nil {
		return fmt.Errorf("failed to close data service connection: %w", err)
	}

	return nil
}
