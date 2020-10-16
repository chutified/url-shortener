package data

import (
	"context"
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
	return nil, nil
}

func (s *service) UpdateRecord(ctx context.Context, id string, r *Record) (*Record, error) {
	//TODO
	return nil, nil
}

func (s *service) DeleteRecord(ctx context.Context, id string) error {
	//TODO
	return nil
}

func (s *service) GetRecordByID(ctx context.Context, id string) (*Record, error) {
	//TODO
	return nil, nil
}

func (s *service) GetRecordByShort(ctx context.Context, short string) (*Record, error) {
	//TODO
	return nil, nil
}

func (s *service) GetRecordByFull(ctx context.Context, full string) (*Record, error) {
	//TODO
	return nil, nil
}

func (s *service) GetAllRecords(ctx context.Context, p int, pagin int) ([]*Record, error) {
	//TODO
	return nil, nil
}
