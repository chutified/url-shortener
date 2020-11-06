package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Config represents the server's settings and the configuration of the database.
type Config struct {
	SrvPort    int    `json:"server_port"`
	SrvTimeOut string `json:"server_timeout"`
	DB         *DB    `json:"db"`
}

// GetConfig returns configuration based on the given file.
// The base of the file's path is at the root of the project (main.go file level).
// Configuration file must be JSON file type (.json).
func GetConfig(file string) (*Config, error) {

	// open config file
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("can not open config file: %w", err)
	}
	defer f.Close()

	// validate file extension
	l := len(file)
	if file[l-5:l] != ".json" {
		return nil, fmt.Errorf("invalid config file type: expected '.json' extension")
	}

	// decode json
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	// validate config's driver
	ok := false
	for _, d := range dialects {
		if cfg.DB.Driver == d {
			ok = true
		}
	}
	if !ok {
		return nil, errors.New("invalid/not suported database dialect: " + cfg.DB.Driver)
	}

	// validate server timeout
	_, err = time.ParseDuration(cfg.SrvTimeOut)
	if err != nil {
		return nil, errors.New("invalid server timeout duration: " + cfg.SrvTimeOut)
	}

	return &cfg, nil
}

// Addr returns a chosen address for the server.
func (cfg *Config) Addr() string {
	return fmt.Sprintf(":%d", cfg.SrvPort)
}