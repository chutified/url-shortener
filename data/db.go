package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/chutified/url-shortener/config"
	"github.com/jmoiron/sqlx"
)

// Service is the controller of the data services.
// Only DataService stores the database connection.
type Service interface {
	InitDB(context.Context, *config.DB) error
	StopDB() error
	AddRecord(context.Context, *Record) (*ShortRecord, error)
	UpdateRecord(context.Context, string, *ShortRecord) (*ShortRecord, error)
	DeleteRecord(context.Context, string) (string, error)
	GetRecordByID(context.Context, string) (*Record, error)
	GetRecordByShort(context.Context, string) (*Record, error)
	GetRecordByShortPeek(context.Context, string) (string, error)
	GetRecordsLen(ctx context.Context) (int, error)
	GetAllRecords(context.Context) ([]*ShortRecord, error)
	RecordRecovery(context.Context, string) (string, error)
}

// service implements Service interface.
type service struct {
	DB *sqlx.DB
}

// NewService is the contructor of the Service controller.
func NewService() Service {
	return &service{}
}

// InitDB intiliazes the database connection for the data server.
// Valid credentials must be provided to connect to the database.
func (s *service) InitDB(ctx context.Context, dbCfg *config.DB) error {

	// retrieve db connection string
	driver, connStr := dbCfg.ConnStr()

	// open connection to db
	var err error
	s.DB, err = sqlx.ConnectContext(ctx, driver, connStr)
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
