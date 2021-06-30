package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chutommy/url-shortener/config"
	"github.com/chutommy/url-shortener/server"
	_ "github.com/lib/pq"
)

func main() {
	// get configuration
	cfg, err := config.GetConfig("settings.json")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config file: %w", err))
	}

	initCtx := context.Background()
	// create server
	srv := server.NewServer()
	err = srv.Set(initCtx, cfg)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to set server: %w", err))
	}

	// run server
	go func() {
		err = srv.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	// wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	log.Printf("Shutting down server... (got signal: %v)\n\n", s)

	// stop server
	if err = srv.Stop(); err != nil {
		log.Println(err.Error())
	}

	// close connections
	if err = srv.Close(); err != nil {
		log.Println(err.Error())
	}
}
