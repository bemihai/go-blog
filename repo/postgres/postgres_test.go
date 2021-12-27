package postgres

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"
)

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "blog"
)

var connection = fmt.Sprintf("postgres://%s:%d/%s?user=%s&password=%s&sslmode=disable", host, port, dbname, user, password)

func TestListArticles(t *testing.T) {

	db, _ := createTestDB(t, connection)
	repo := PSQLRepository{DB: db}

	// table containing 2 entries
	dumpTestData(t, db)
	articles, err := repo.ListArticles()
	require.NotEmpty(t, articles)
	require.NoError(t, err)
	require.Len(t, articles, 2)

	// empty table
	truncateTables(t, db)
	articles, err = repo.ListArticles()
	require.Empty(t, articles)
	require.NoError(t, err)
	require.Len(t, articles, 0)
}

func TestListAuthors(t *testing.T) {

	db, _ := createTestDB(t, connection)
	repo := PSQLRepository{DB: db}

	// table containing 2 entries
	dumpTestData(t, db)
	authors, err := repo.ListAuthors()
	require.NotEmpty(t, authors)
	require.NoError(t, err)
	require.Len(t, authors, 2)

	// empty table
	truncateTables(t, db)
	authors, err = repo.ListAuthors()
	require.Empty(t, authors)
	require.NoError(t, err)
	require.Len(t, authors, 0)
}

func TestGetArticleById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	repo := PSQLRepository{DB: db}

	dumpTestData(t, db)

	// existing article
	a, err := repo.GetArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403126")
	require.NotEmpty(t, a)
	require.NoError(t, err)
	require.Equal(t, a.Title, "Test title 1")

	// invalid uuid
	_, err = repo.GetArticleById("invalid uuid")
	require.Error(t, err)

	// non-existing article
	_, err = repo.GetArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
	require.ErrorIs(t, err, ErrArticleNotFound)
}

func TestGetAuthorById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	repo := PSQLRepository{DB: db}

	dumpTestData(t, db)

	// existing author
	a, err := repo.GetAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403124")
	require.NotEmpty(t, a)
	require.NoError(t, err)
	require.Equal(t, a.Name, "Test Author1")

	// invalid uuid
	_, err = repo.GetAuthorById("invalid uuid")
	require.Error(t, err)

	// non-existing author
	_, err = repo.GetAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
	require.ErrorIs(t, err, ErrAuthorNotFound)
}
