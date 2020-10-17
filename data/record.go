package data

import (
	"context"
	"fmt"
	"time"
)

// Record is the unit of each shorten URL.  Record stores the time of its creation,
// uppdate and deletion. All Short atributes must be unique. Full can have duplicates.
type Record struct {
	ID        string
	Full      string
	Short     string
	Usage     int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

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

func (s *service) DeleteRecord(ctx context.Context, id string) error {
	//TODO
	fmt.Printf("Deleted record id (%s)\n", id)
	return nil
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

func (s *service) GetAllRecords(ctx context.Context, p int, pagin int) ([]*Record, error) {
	//TODO
	fmt.Printf("Served records (page %d, with pagin %d)\n", p, pagin)
	return nil, nil
}
