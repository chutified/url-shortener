package data

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Record is the unit of each shorten URL.  Record stores the time of its creation,
// update and deletion. All Short atributes must be unique. Full can have duplicates.
type Record struct {
	ID        string    `json:"id"`
	Full      string    `json:"full"`
	Short     string    `json:"short"`
	Usage     int32     `json:"usage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

var (
	// ErrInvalidRecord is returned when an invalid record is provided.
	ErrInvalidRecord = errors.New("given record is invalid")
	// ErrShortNotFound is returned when record's ID can not be found.
	ErrIDNotFound = errors.New("given 'id' does not exist")
	// ErrShortNotFound is returned when record's Short can not be found.
	ErrShortNotFound = errors.New("given 'short' does not exist")
	// ErrFullNotFound is returned when record's Full can not be found.
	ErrFullNotFound = errors.New("given 'full' does not exist")
	// ErrUnavailableShort is returned when the new record has a Short which already exists.
	ErrUnavailableShort = errors.New("short of given record is already in use")
)

func (s *service) AddRecord(ctx context.Context, r *Record) (*Record, error) {
	//TODO
	fmt.Println("Added record:", *r)
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

func (s *service) GetAllRecords(ctx context.Context, sort string, p int, pagin int) ([]*Record, error) {
	//TODO
	fmt.Printf("Served records sorted by %s (page %d, with pagin %d)\n", sort, p, pagin)
	return nil, nil
}
