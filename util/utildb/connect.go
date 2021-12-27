package utildb

import "database/sql"

// Connect connects to a database and verifies the connection with a ping.
func Connect(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close() // nolint: errcheck
		return nil, err
	}
	return db, nil
}

// MustConnect connects to a database and verifies the connection with a ping.
// Panics on error.
func MustConnect(driverName, dataSourceName string) *sql.DB {
	db, err := Connect(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}
