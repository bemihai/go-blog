package postgres

import (
	repo "blog/repo"
	"blog/util/utildb"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// test data
var articles = []repo.Article{
	{
		Id:       "b4a4de9e-2f52-4cf1-8907-3d828d403126",
		Title:    "Test title 1",
		Body:     "Test body 1",
		PostedAt: time.Now(),
		Author: repo.Author{
			Id:    "b4a4de9e-2f52-4cf1-8907-3d828d403124",
			Name:  "Test Author1",
			Email: "test.author1@email.com",
		},
	},
	{
		Id:       "b4a4de9e-2f52-4cf1-8907-3d828d403127",
		Title:    "Test title 2",
		Body:     "Test body 2",
		PostedAt: time.Now(),
		Author: repo.Author{
			Id:    "b4a4de9e-2f52-4cf1-8907-3d828d403125",
			Name:  "Test Author2",
			Email: "test.author2@email.com",
		},
	},
}

// createTestDB sets up a random schema in the given db.
// It returns the connection to the database and the schema name. Schema is
// dropped and DB is closed when the test is finished.
func createTestDB(tb testing.TB, connection string) (*sql.DB, string) {
	tb.Helper()

	db, schema, err := utildb.SwitchToRandomSchema(connection)
	require.NoError(tb, err, "Could not create blog schema")
	require.NotEmpty(tb, schema, "Empty schema")

	err = utildb.ExecFile(db, "../../db/blog_migrate.sql")
	require.NoError(tb, err, "Could not create tables")

	err = dumpTestData(db)
	if err != nil {
		panic(err)
	}

	tb.Cleanup(func() {
		utildb.DropSchema(db, schema) // nolint: errcheck
		db.Close()                    // nolint: errcheck
	})

	return db, schema
}

// truncateTables truncates tables in the given db.
func truncateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM articles; DELETE FROM authors;")
	require.NoError(t, err, "Could not truncate tables")
}

// dumpTestData dumps the test data to the given db.
func dumpTestData(db *sql.DB) error {

	for _, art := range articles {

		query := `INSERT INTO authors(id, name, email) values ($1, $2, $3);`
		_, err := db.Exec(query, art.Author.Id, art.Author.Name, art.Author.Email)
		if err != nil {
			return fmt.Errorf("could not add authors: %w", err)
		}

		query = `INSERT INTO articles(id, title, body, author_id) values ($1, $2, $3, $4);`
		_, err = db.Exec(query, art.Id, art.Title, art.Body, art.Author.Id)
		if err != nil {
			return fmt.Errorf("could not add articles: %w", err)
		}
	}
	return nil
}
