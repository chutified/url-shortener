package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	ErrInvalidRecord = errors.New("given record is invalid")
	// ErrIDNotFound is returned when record's ID can not be found.
	ErrIDNotFound = errors.New("given 'id' does not exist")
	// ErrShortNotFound is returned when record's Short can not be found.
	ErrShortNotFound = errors.New("given 'short' does not exist")
	// ErrFullNotFound is returned when record's Full can not be found.
	ErrFullNotFound = errors.New("given 'full' does not exist")
	// ErrUnavailableShort is returned when the new record has a Short which already exists.
	ErrUnavailableShort = errors.New("short of given record is already in use")
)

func (s *service) AddRecord(ctx context.Context, r *Record) (*Record, error) {

	// validate values
	if r.Full == "" || r.Short == "" {
		return nil, ErrInvalidRecord
	}

	// process record
	r = &Record{
		ID:    uuid.New().String(),
		Full:  strings.ToLower(r.Full),
		Short: strings.ToLower(r.Short),
		Usage: 0,
	}

	return nil, nil
}

func (s *service) UpdateRecord(ctx context.Context, id string, r *Record) (*Record, error) {
	//TODO
	fmt.Printf("Updated record (id %s): %v\n", id, *r)
	return nil, nil
}

func (s *service) DeleteRecord(ctx context.Context, id string) (int, error) {
	//TODO
	fmt.Printf("Deleted record id (%s)\n", id)
	return 0, nil
}

func (s *service) GetRecordByID(ctx context.Context, id string) (*Record, error) {
	//TODO
	fmt.Printf("Served record (id %s)\n", id)
	return nil, nil
}

func (s *service) GetRecordByShort(ctx context.Context, short string) (*Record, error) {
	//TODO
	fmt.Printf("Served record (short %s)\n", short)
	return nil, nil
}

func (s *service) GetRecordByFull(ctx context.Context, full string) (*Record, error) {
	//TODO
	fmt.Printf("Served record (full %s)\n", full)
	return nil, nil
}

func (s *service) GetRecordsLen(ctx context.Context) (int, error) {
	//TODO
	fmt.Println("Served number of records.")
	return 0, nil
}

func (s *service) GetAllRecords(ctx context.Context, pcfg PageCfg) ([]*Record, PageCfg, error) {
	//TODO
	fmt.Printf("Served records sorted by %s (page %d, with pagin %d)\n", pcfg.Sort, pcfg.Page, pcfg.Pagin)
	return nil, pcfg, nil
}

func (s *service) GetRecordDetails(ctx context.Context, id string) (*Record, error) {
	//TODO
	fmt.Printf("Served record's details (id %s)\n", id)
	return nil, nil
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
