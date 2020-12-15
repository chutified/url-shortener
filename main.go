package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chutified/url-shortener/config"
	"github.com/chutified/url-shortener/controller"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {

	// get configuration
	cfg, err := config.GetConfig("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	initCtx := context.Background()
	// create server
	srv := controller.NewServer()
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
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// stop server
	if err = srv.Stop(); err != nil {
		log.Println(err.Error())
	}

	// close connections
	if err = srv.Close(); err != nil {
		log.Println(err.Error())
	}
}
