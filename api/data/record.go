package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Record is the unit of each shorten URL.  Record stores the time of its creation,
// update and deletion. All Short atributes must be unique. Full can have duplicates.
type Record struct {
	ID        string       `json:"shortcut_id"`
	Full      string       `json:"full_url"`
	Short     string       `json:"short_url"`
	Usage     int32        `json:"usage"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// PageCfg holds a configuration for retrieving large number of records.
type PageCfg struct {
	Page  int    `json:"page"`
	Pagin int    `json:"pagin"`
	Sort  string `json:"sort"`
}

func (p *PageCfg) checkSort() {

	// check if Sort attribute is valid
	sortable := []string{"shortcut_id", "full", "short", "usage", "created_at", "updated_at"}
	for _, s := range sortable {
		if p.Sort == s {
			return
		}
	}

	// set default Sort value
	p.Sort = "shortcut_id"
}

var (
	// ErrInvalidRecord is returned when an invalid record is provided.
	ErrInvalidRecord = errors.New("given record is unavailable")
	// ErrIDNotFound is returned when record's ID can not be found.
	ErrIDNotFound = errors.New("given 'id' value does not exist")
	// ErrShortNotFound is returned when record's Short can not be found.
	ErrShortNotFound = errors.New("given 'short' value does not exist")
	// ErrFullNotFound is returned when record's Full can not be found.
	ErrFullNotFound = errors.New("given 'full' value does not exist")
	// ErrUnavailableShort is returned when the new record has a Short which already exists.
	ErrUnavailableShort = errors.New("'short' value of the given record is already in use")
	// ErrInvalidID is returned when an ID with invalid formatis provided.
	ErrInvalidID = errors.New("given id has invalid format")
)

// AddRecord inserts a new record into the database. Only Full and Short
// record's attributes must be set, other are omitted.
// If any error occurs ErrInvalidRecord, ErrUnavailableShort or an unexpected
// internal server error is returned.
func (s *service) AddRecord(ctx context.Context, r *Record) (*Record, error) {

	// validate values
	if r.Full == "" || r.Short == "" {
		return nil, ErrInvalidRecord
	}

	// create a record
	newr := &Record{
		ID:    uuid.New().String(),
		Full:  strings.ToLower(r.Full),
		Short: strings.ToLower(r.Short),
	}

	// insert record
	_, err := s.DB.ExecContext(ctx, `
INSERT INTO
  shortcuts (shortcut_id, full_url, short_url)
VALUES
  ($1, $2, $3);
  `, newr.ID, newr.Full, newr.Short)
	if err != nil {

		// postgres errors
		if err, ok := err.(*pq.Error); ok {
			// unique violation
			if err.Code == "23505" {
				return nil, ErrUnavailableShort
			}
		}

		return nil, fmt.Errorf("could not execute sql insert: %w", err)
	}

	return newr, nil
}

// UpdateRecord updates a record with the given id.
// If the record is not found ErrIDNotFound is returned.
// If the record r has a short which is already in use,
// ErrUnavailableShort is returned. Any other errors
// are server internal.
func (s *service) UpdateRecord(ctx context.Context, id string, r *Record) (*Record, error) {

	// create record
	updr := &Record{
		ID:    strings.ToLower(id),
		Full:  strings.ToLower(r.Full),
		Short: strings.ToLower(r.Short),
	}

	// update record
	result, err := s.DB.ExecContext(ctx, `
UPDATE
  shortcuts
SET
  full_url = COALESCE($2, full_url),
  short_url = COALESCE($3, short_url)
WHERE
  shortcut_id = $1
  AND deleted_at IS NULL;
	`, updr.ID, newNullString(updr.Full), newNullString(updr.Short))
	if err != nil {

		// postgres errors
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {

			// unique violation
			case "23505":
				return nil, ErrUnavailableShort

			// invalid text representation
			case "22P02":
				return nil, ErrInvalidID
			}
		}

		return nil, fmt.Errorf("could not execute sql update; %w", err)
	}

	// check affected row
	if n, _ := result.RowsAffected(); n != 1 {
		return nil, ErrIDNotFound
	}

	return updr, nil
}

// DeleteRecord softly removes a record with the given id.
// On success the function returns an ID of the deleted record
// as a non-empty string.
func (s *service) DeleteRecord(ctx context.Context, id string) (string, error) {

	// id to lowercase
	id = strings.ToLower(id)

	// softly remove record
	result, err := s.DB.ExecContext(ctx, `
UPDATE
  shortcuts
SET
  deleted_at = NOW()
WHERE
  shortcut_id = $1
  AND deleted_at IS NULL;
	`, id)
	if err != nil {

		// postres errors
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {

			// invalid text representation
			case "22P02":
				return "", ErrInvalidID
			}
		}

		return "", fmt.Errorf("could not execute sql update; %w", err)
	}

	// check result
	if n, _ := result.RowsAffected(); n != 1 {
		return "", ErrIDNotFound
	}

	return id, nil
}

// GetRecordByID find a record in the database by uuid.
// If no row is found the function returns ErrIDNotFound or
// in case the given id does not follow standard UUID format, ErrInvalidID instead,
// If any unexpected error occurs, unexpected server error is returned.
func (s *service) GetRecordByID(ctx context.Context, id string) (*Record, error) {

	// id to lowercase
	id = strings.ToLower(id)

	// query records
	row := s.DB.QueryRowContext(ctx, `
SELECT
  shortcut_id,
  full_url,
  short_url,
  usage,
  created_at,
  updated_at
FROM
  shortcuts
WHERE
  shortcut_id = $1
  AND deleted_at IS NULL
LIMIT 1;
	`, id)

	// scan row into new record
	var r Record
	err := row.Scan(&r.ID, &r.Full, &r.Short, &r.Usage, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		// nothing returned
		return nil, ErrIDNotFound

	} else if err != nil {

		// postgres errors
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {

			// invalid text representation
			case "22P02":
				return nil, ErrInvalidID
			}
		}

		return nil, fmt.Errorf("unexpected query error: %w", err)
	}

	return &r, nil
}

// GetRecordByShort is an alternative of GetRecordByID which uses
// short atribute for querying instead of an ID.
func (s *service) GetRecordByShort(ctx context.Context, short string) (*Record, error) {

	// short to lowercase
	short = strings.ToLower(short)

	// select the record
	row := s.DB.QueryRowContext(ctx, `
SELECT
  shortcut_id,
  full_url,
  short_url,
  usage,
  created_at,
  updated_at
FROM
  shortcuts
WHERE
  short_url = $1
  AND deleted_at IS NULL
LIMIT 1;
	`, short)

	// scan row into a new record
	var r Record
	err := row.Scan(&r.ID, &r.Full, &r.Short, &r.Usage, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		// nothing resturned
		return nil, ErrShortNotFound

	} else if err != nil {
		return nil, fmt.Errorf("unexpected sql query error: %w", err)
	}

	return &r, nil
}

// GetRecordsLen returns the number of active urls.
// Only unexpected errors can occurre.
func (s *service) GetRecordsLen(ctx context.Context) (int, error) {

	// query database
	row := s.DB.QueryRowContext(ctx, `
SELECT
  COUNT(*)
FROM
  shortcuts
WHERE
  deleted_at IS NOT NULL;
	`)

	// scan row
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unexpected sql query error: %w", err)
	}

	return count, nil
}

func (s *service) GetAllRecords(ctx context.Context, pcfg PageCfg) ([]*Record, PageCfg, error) {
	//TODO
	fmt.Printf("Served records sorted by %s (page %d, with pagin %d)\n", pcfg.Sort, pcfg.Page, pcfg.Pagin)
	return nil, pcfg, nil
}

func (s *service) RecordRecovery(ctx context.Context, id string) (*Record, error) {
	//TODO
	fmt.Printf("Deleted record recoverd (id %s)\n", id)
	return nil, nil
}

func (s *service) incrementUsage(ctx context.Context, id string) (*Record, error) {
	// TODO
	fmt.Printf("Usage incremented (id %s)\n", id)
	return nil, nil
}

func (s *service) TotalUsage(ctx context.Context) (int, error) {
	// TODO
	fmt.Println("Total usage served")
	return -1, nil
}
