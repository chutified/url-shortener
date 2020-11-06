package controller

import (
	"context"
	"fmt"

	"github.com/chutified/url-shortener/api/config"
	"github.com/chutified/url-shortener/api/data"
)

// handler is the controller of the data service actions.
type handler struct {
	ds data.Service
}

// newHandler returns an empty handler.
func newHandler() *handler {
	return &handler{}
}

// initDataService initializes handler's data service.
func (h *handler) initDataService(ctx context.Context, dbCfg *config.DB) error {

	// create new data service
	h.ds = data.NewService()

	// initialize data service
	err := h.ds.InitDB(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize data service: %w", err)
	}

	return nil
}

// closeHandler stops all active connections. closeHandler closes the data service.
// This function should not be called often (meant to be used only when the server
// is shutting down).
func (h *handler) closeHandler() error {

	// close data service connection
	err := h.ds.StopDB()
	if err != nil {
		return fmt.Errorf("failed to close data service connection")
	}

	return nil
}
