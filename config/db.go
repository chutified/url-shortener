package config

// DB holds credentials of the database.
type DB struct {
	Driver string `json:"driver"`
	DBConn string `json:"-"`
}

// ConnStr returns a driver and a connection string of the database.
func (db *DB) ConnStr() (string, string) {
	// return "postgres", fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
	// db.Host, db.Port, db.User, db.Password, db.DBName)
	return "postgres", db.DBConn
}
