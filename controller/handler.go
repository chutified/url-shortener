package controller

import (
	"fmt"

	"github.com/chutified/url-shortener/config"
	"github.com/chutified/url-shortener/data"
)

// handler is the controller of the data service actions.
type handler struct {
	ds data.Service
}

// initHandler create a new handler with a data service.
func (h *handler) initHandler(dbCfg *config.DB) (*handler, error) {

	// create new data service
	h.ds = data.NewService()

	// initialize data service
	err := h.ds.InitDB(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize data service: %w", err)
	}

	return h, nil
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
