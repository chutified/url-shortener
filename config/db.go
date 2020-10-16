package config

import "fmt"

// DB holds credentials of the Postgres database.
type DB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// ConnStr returns a connection string
// for the psql database.
func (db *DB) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DBName)
}
