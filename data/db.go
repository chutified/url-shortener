package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chutified/url-shortener/config"
)

// Service is the controller of the data services.
// Only DataService stores the database connection.
type Service struct {
	DB *sql.DB
}

// NewService is the contructor of the Service controller.
func NewService() *Service {
	return &Service{}
}

// InitDB intiliazes the database connection for the data server.
// Valid credentials must be provided to connect to the database.
func (s *Service) InitDB(dbCfg config.DB) error {

	// retrieve db connection string
	connStr := dbCfg.ConnStr()

	// open connection to db
	var err error
	s.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open db conn: %w", err)
	}

	// test connection
	err = s.DB.Ping()
	if err != nil {
		return errors.New("db conn verification failed")
	}

	return nil
}
