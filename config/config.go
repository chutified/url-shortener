package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	// ErrFileNotFound is returned if file with given address can not be found.
	ErrFileNotFound = errors.New("can not open config file")
	// ErrInvalidFileFormat is returned if configuration file with unsupported file extension is provided.
	ErrInvalidFileFormat = errors.New("invalid config file type: expected '.json' extension")
	// ErrInvalidJSONFile is returned if config file's content is corrupted.
	ErrInvalidJSONFile = errors.New("unable to correctly decode config file")
	// ErrDriverNotSupported is returned if unsupported sql driver is provided.
	ErrDriverNotSupported = errors.New("only postgres sql driver is supported")
	// ErrInvalidTimeFormat is returned if given time can not be correctly formatted.
	ErrInvalidTimeFormat = errors.New("invalid server timeout duration")
	// ErrDBCONNEnvVarNotSet is returned if environment variable of the database connection is not set.
	ErrDBCONNEnvVarNotSet = errors.New(
		"environment variable of url (URL_SHORTENER_DBCONN) for database connection is not set")
)

// Config represents the server's settings and the configuration of the database.
type Config struct {
	SrvPort    int    `json:"server_port"`
	SrvTimeOut string `json:"server_timeout"`
	DB         *DB    `json:"db"`
}

// GtConfig returns configuration based on the given file.
// The base of the file's path is at the root of the project (main.go file level).
// Configuration file must be JSON file type (.json).
func GetConfig(file string) (*Config, error) {
	// get config file
	cfg, err := OpenConfig(file)
	if err != nil {
		return nil, fmt.Errorf("could not correctly load config file: %w", err)
	}

	// check for database connection environment variable
	dbConn := os.Getenv("URL_SHORTENER_DBCONN")
	if dbConn == "" {
		return nil, ErrDBCONNEnvVarNotSet
	}

	cfg.DB.DBConn = dbConn

	return &cfg, nil
}

// Addr returns a chosen address for the server.
func (cfg *Config) Addr() string {
	if cfg != nil {
		return fmt.Sprintf(":%d", cfg.SrvPort)
	}

	return ""
}

// OpenConfig load config file, validate its extension and decode it into a Config struct.
func OpenConfig(file string) (Config, error) { // open config file
	f, err := os.Open(file)
	if err != nil {
		return Config{}, ErrFileNotFound
	}
	defer f.Close()

	// validate file extension
	l := len(file)

	if file[l-5:l] != ".json" {
		return Config{}, ErrInvalidFileFormat
	}

	// decode json
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)

	if err != nil {
		return Config{}, ErrInvalidJSONFile
	}

	// validate config's driver
	if cfg.DB.Driver != "postgres" {
		return Config{}, ErrDriverNotSupported
	}

	// validate server timeout
	_, err = time.ParseDuration(cfg.SrvTimeOut)
	if err != nil {
		return Config{}, ErrInvalidTimeFormat
	}

	return cfg, nil
}
