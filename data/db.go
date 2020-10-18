package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chutified/url-shortener/config"
)

// Service is the controller of the data services.
// Only DataService stores the database connection.
type Service interface {
	InitDB(*config.DB) error
	StopDB() error
	AddRecord(context.Context, *Record) (*Record, error)
	UpdateRecord(context.Context, string, *Record) (*Record, error)
	DeleteRecord(context.Context, string) error
	GetRecordByID(context.Context, string) (*Record, error)
	GetRecordByShort(context.Context, string) (*Record, error)
	GetRecordByFull(context.Context, string) (*Record, error)
	GetAllRecords(context.Context, string, int, int) ([]*Record, error)
}

// service implements Service interface.
type service struct {
	DB *sql.DB
}

// NewService is the contructor of the Service controller.
func NewService() Service {
	return &service{}
}

// InitDB intiliazes the database connection for the data server.
// Valid credentials must be provided to connect to the database.
func (s *service) InitDB(dbCfg *config.DB) error {

	// retrieve db connection string
	driver, connStr := dbCfg.ConnStr()

	// open connection to db
	var err error
	s.DB, err = sql.Open(driver, connStr)
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

// StopDB closes database connection of the service.
func (s *service) StopDB() error {

	// close db conection
	err := s.DB.Close()
	if err != nil {
		return err
	}

	return nil
}
