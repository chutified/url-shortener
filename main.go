package main

import (
	"fmt"
	"log"

	"github.com/chutified/url-shortener/config"
	"github.com/chutified/url-shortener/data"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cfg, err := config.GetConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	s := data.NewService()
	err = s.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer s.StopDB()

	fmt.Println("OK")
}
