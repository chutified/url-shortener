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

// ErrPQInvalidTextRepresentation occurs when type is represented in a bad format.
const ErrPQInvalidTextRepresentation = "22P02"

// ErrPQUniqueKeyViolation occurs when Postgres' unique key is violated.
const ErrPQUniqueKeyViolation = "23505"

// Record is the unit of each shorten URL.  Record stores the time of its creation,
// update and deletion. All Short attributes must be unique. Full can have duplicates.
type Record struct {
	ID        string       `json:"shortcut_id"`
	Full      string       `json:"full_url"`
	Short     string       `json:"short_url"`
	Usage     int32        `json:"usage"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

// ShortRecord represents shorter version of the Record.
type ShortRecord struct {
	ID    string `json:"shortcut_id"`
	Full  string `json:"full_url"`
	Short string `json:"short_url"`
	Usage int32  `json:"usage"`
}

var (
	// ErrInvalidRecord is returned when an invalid record is provided.
	ErrInvalidRecord = errors.New("given record is invalid ")
	// ErrIDNotFound is returned when record's ID can not be found.
	ErrIDNotFound = errors.New("given 'id' value does not exist")
	// ErrShortNotFound is returned when record's Short can not be found.
	ErrShortNotFound = errors.New("given 'short' value does not exist")
	// ErrUnavailableShort is returned when the new record has a Short which already exists.
	ErrUnavailableShort = errors.New("'short' value of the given record was already used")
	// ErrInvalidID is returned when an ID with invalid formats provided.
	ErrInvalidID = errors.New("given id has invalid format")
	// ErrNotDeleted is returned when the record is not deleted.
	ErrNotDeleted = errors.New("record with the given id is either not deleted or does not exist")
)

// AddRecord inserts a new record into the database. Only Full and Short
// record's attributes must be set, other are omitted.
// If any error occurs ErrInvalidRecord, ErrUnavailableShort or an unexpected
// internal server error is returned.
func (s *service) AddRecord(ctx context.Context, r *Record) (*ShortRecord, error) {
	// validate values
	if r.Full == "" || r.Short == "" {
		return nil, ErrInvalidRecord
	}

	// create a record
	newRec := &ShortRecord{
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
  `, newRec.ID, newRec.Full, newRec.Short)
	if err != nil {
		// postgres errors
		var pqErr *pq.Error
		if errors.As(err, pqErr) {
			// unique violation
			if pqErr.Code == ErrPQUniqueKeyViolation {
				return nil, ErrUnavailableShort
			}
		}

		return nil, fmt.Errorf("could not execute sql insert: %w", err)
	}

	return newRec, nil
}

// UpdateRecord updates a record with the given id.
// If the record is not found ErrIDNotFound is returned.
// If the record r has a short which is already in use,
// ErrUnavailableShort is returned. Any other errors
// are server internal.
func (s *service) UpdateRecord(ctx context.Context, id string, r *ShortRecord) (*ShortRecord, error) {
	// create record
	updRecord := &ShortRecord{
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
  `, updRecord.ID, newNullString(updRecord.Full), newNullString(updRecord.Short))
	if err != nil {
		// postgres errors
		var pqErr *pq.Error
		if errors.As(err, pqErr) {
			switch pqErr.Code {
			// unique violation
			case ErrPQUniqueKeyViolation:
				return nil, ErrUnavailableShort

			// invalid text representation
			case ErrPQInvalidTextRepresentation:
				return nil, ErrInvalidID
			}
		}

		return nil, fmt.Errorf("could not execute sql update; %w", err)
	}

	// check affected row
	if n, _ := result.RowsAffected(); n != 1 {
		return nil, ErrIDNotFound
	}

	return updRecord, nil
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
		// postgres errors
		// if err, ok := err.(*pq.Error); ok {
		var pqErr *pq.Error
		if errors.As(err, pqErr) {
			if pqErr.Code == ErrPQInvalidTextRepresentation {
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

	if errors.Is(err, sql.ErrNoRows) {
		// nothing returned
		return nil, ErrIDNotFound
	} else if err != nil {
		// postgres errors
		var pqErr *pq.Error
		if errors.As(err, pqErr) {
			if pqErr.Code == ErrPQInvalidTextRepresentation {
				return nil, ErrInvalidID
			}
		}

		return nil, fmt.Errorf("unexpected query error: %w", err)
	}

	return &r, nil
}

// GetRecordByShort is an alternative of GetRecordByID which uses
// short attribute for querying instead of an ID.
func (s *service) GetRecordByShort(ctx context.Context, short string) (*Record, error) { // short to lowercase
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

	if errors.Is(err, sql.ErrNoRows) {
		// nothing returned
		return nil, ErrShortNotFound
	} else if err != nil {
		return nil, fmt.Errorf("unexpected sql query error: %w", err)
	}

	return &r, nil
}

// GetRecordByShortPeek finds and returns the full version which corresponds
// to the given short url.
func (s *service) GetRecordByShortPeek(ctx context.Context, short string) (string, error) { // short to lowercase
	short = strings.ToLower(short)

	// get full url
	row := s.DB.QueryRowContext(ctx, `
SELECT
  shortcut_id,
  full_url
FROM
  shortcuts
WHERE
  short_url = $1
  AND deleted_at IS NULL
LIMIT 1;
  `, short)

	// scan full url
	var id, full string
	err := row.Scan(&id, &full)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrShortNotFound
	} else if err != nil {
		return "", fmt.Errorf("unexpected sql query error: %w", err)
	}
	// increment record's usage
	if err = s.incrementUsage(ctx, id); err != nil {
		return "", err
	}

	// log record's usage
	if err = s.logUsage(ctx, id); err != nil {
		return "", err
	}

	return full, nil
}

// GetRecordsLen returns the number of active urls.
// Only unexpected errors can occur.
func (s *service) GetRecordsLen(ctx context.Context) (int, error) {
	// query database
	row := s.DB.QueryRowContext(ctx, `
SELECT
  COUNT(*) as C
FROM
  shortcuts
WHERE
  deleted_at IS NULL;
  `)

	// scan row
	var count int

	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unexpected sql query error: %w", err)
	}

	return count, nil
}

// GetAllRecords returns all active records in the database.
func (s *service) GetAllRecords(ctx context.Context) (records []*ShortRecord, err error) {
	// retrieve all records
	rows, _ := s.DB.QueryContext(ctx, `
SELECT
  shortcut_id,
  full_url,
  short_url,
  usage
FROM
  shortcuts
WHERE
  deleted_at is NULL
ORDER BY
  usage;
  `)
	// if err != nil {
	// 	return nil, fmt.Errorf("unexpected server error while scanning racords: %w", err)
	// }

	defer func() {
		err = rows.Close()
	}()

	for rows.Next() {
		// create new record
		var r ShortRecord

		if err := rows.Scan(&r.ID, &r.Full, &r.Short, &r.Usage); err != nil {
			return nil, fmt.Errorf("unexpected server error while scanning racords: %w", err)
		}

		// store
		records = append(records, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected server error while scanning racords: %w", err)
	}

	return records, nil
}

// RecordRecovery recovers the softly deleted records.
func (s *service) RecordRecovery(ctx context.Context, id string) (string, error) {
	// id to lowercase
	id = strings.ToLower(id)

	// recover removed Record
	result, err := s.DB.ExecContext(ctx, `
UPDATE
  shortcuts
SET
  deleted_at = NULL
WHERE
  shortcut_id = $1
  AND deleted_at IS NOT NULL;
  `, id)
	if err != nil {
		// postgres errors
		var pqErr *pq.Error
		if errors.As(err, pqErr) {
			if pqErr.Code == ErrPQInvalidTextRepresentation {
				return "", ErrInvalidID
			}
		}

		return "", fmt.Errorf("could not execute sql update; %w", err)
	}

	// check result
	if n, _ := result.RowsAffected(); n != 1 {
		return "", ErrNotDeleted
	}

	return id, nil
}
