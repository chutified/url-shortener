package data

import "time"

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
