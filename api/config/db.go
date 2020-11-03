package config

import "fmt"

// DB holds credentials of the database.
type DB struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var dialects = []string{"mysql", "postgres"}

// ConnStr returns a driver and a connection string of the database.
func (db *DB) ConnStr() (string, string) {

	switch db.Driver {

	case "postgres":
		return "postgres", fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.DBName)

	case "mysql":
		return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			db.User, db.Password, db.Host, db.Port, db.DBName)

	default:
		return "", ""
	}
}
