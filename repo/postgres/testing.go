package postgres

import (
	repo "blog/repo"
	"blog/util/utildb"
	"blog/util/utiltesting"
	"database/sql"
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

	tb.Cleanup(func() {
		utildb.DropSchema(db, schema) // nolint: errcheck
		db.Close()                    // nolint: errcheck
	})

	err = utildb.ExecFile(db, utiltesting.AbsolutePath("/blog/sql/test.sql"))
	require.NoError(tb, err, "Could not create tables")

	return db, schema
}

// truncateTables truncates tables in the given db.
func truncateTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM authors; DELETE FROM articles;")
	require.NoError(t, err, "Could not truncate tables")
}

// dumpTestData dumps the test data to the given db.
func dumpTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, art := range articles {
		query := `INSERT INTO authors(id, name, email) values ($1, $2, $3)`
		_, err := db.Exec(query, art.Author.Id, art.Author.Name, art.Author.Email)
		require.NoError(t, err, "Could not add authors")
		query = `INSERT INTO articles(id, title, body, author_id) values ($1, $2, $3, $4)`
		_, err = db.Exec(query, art.Id, art.Title, art.Body, art.Author.Id)
		require.NoError(t, err, "Could not add articles")
	}
}
