package config

// DB holds credentials of the database.
type DB struct {
	Driver string `json:"driver"`
	DBConn string `json:"-"`
}

// ConnStr returns a driver and a connection string of the database.
func (db *DB) ConnStr() (string, string) {
	return "postgres", db.DBConn
}
