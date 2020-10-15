package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chutified/url-shortener/config"
)

// InitDB intiliazes the database connection.
// Valid credentials must be provided to connect the
// Postgres database.
func InitDB(dbCfg config.DB) (*sql.DB, error) {

	// open connection to db
	db, err := sql.Open("postgres", dbCfg.ConnStr())
	if err != nil {
		return nil, fmt.Errorf("failed to open db conn: %w", err)
	}

	// test connection
	err = db.Ping()
	if err != nil {
		return nil, errors.New("db conn verification failed")
	}

	return db, nil
}
