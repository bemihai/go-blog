package utildb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// SwitchToSchema creates a new schema on the given database,
// and returns a database handle with search_path reflecting the new schema.
// The connection string must be a URL (e.g: psql "postgres://localhost:5434/postgres?sslmode=disable").
func SwitchToSchema(conn string, schema string) (*sql.DB, error) {
	db, err := Connect("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err = CreateSchema(db, schema); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}
	if err = db.Close(); err != nil {
		return nil, fmt.Errorf("close db: %w", err)
	}

	conn = fmt.Sprintf("%s&search_path=%s,public", conn, schema)

	db, err = Connect("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return db, nil
}

// SwitchToRandomSchema creates a new schema with a random name on the given database,
// and returns a database handle with search_path reflecting the new schema.
// The connection string must be a URL (e.g: psql "postgres://localhost:5434/postgres?sslmode=disable").
func SwitchToRandomSchema(conn string) (*sql.DB, string, error) {
	schema := fmt.Sprintf("schema_%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	db, err := SwitchToSchema(conn, schema)
	if err != nil {
		return nil, "", err
	}
	return db, schema, nil
}

// CreateRandomSchema creates a new schema with a random name on the given database.
func CreateRandomSchema(db *sql.DB) (string, error) {
	schema := fmt.Sprintf("schema_%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	if err := CreateSchema(db, schema); err != nil {
		return "", err
	}
	return schema, nil
}

// CreateSchema creates a new schema on the given database.
func CreateSchema(db *sql.DB, schema string) error {
	if err := DropSchemaIfExists(db, schema); err != nil {
		return fmt.Errorf("drop schema: %w", err)
	}

	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema))
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}

// DropSchema drops the given schema from the db.
func DropSchema(db *sql.DB, schema string) error {
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", schema))
	return err
}

// DropSchemaIfExists drops the given schema, if it exists, from the db.
func DropSchemaIfExists(db *sql.DB, schema string) error {
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
	return err
}
