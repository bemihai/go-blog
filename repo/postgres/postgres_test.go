package postgres

import (
	repo "blog/repo"
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "blog"
)

var connection = fmt.Sprintf("postgres://%s:%d/%s?user=%s&password=%s&sslmode=disable", host, port, dbname, user, password)

func TestListArticles(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	// table containing 2 entries
	dumpTestData(t, db)
	articles, err := r.ListArticles()
	require.NotEmpty(t, articles)
	require.NoError(t, err)
	require.Len(t, articles, 2)

	// empty table
	truncateTables(t, db)
	articles, err = r.ListArticles()
	require.Empty(t, articles)
	require.NoError(t, err)
	require.Len(t, articles, 0)

	// closed connection
	db.Close()
	_, err = r.ListArticles()
	require.Error(t, err)
}

func TestListAuthors(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	// table containing 2 entries
	dumpTestData(t, db)
	authors, err := r.ListAuthors()
	require.NotEmpty(t, authors)
	require.NoError(t, err)
	require.Len(t, authors, 2)

	// empty table
	truncateTables(t, db)
	authors, err = r.ListAuthors()
	require.Empty(t, authors)
	require.NoError(t, err)
	require.Len(t, authors, 0)

	// closed connection
	db.Close()
	_, err = r.ListAuthors()
	require.Error(t, err)

}

func TestGetArticleById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	dumpTestData(t, db)

	// existing article
	a, err := r.GetArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403126")
	require.NotEmpty(t, a)
	require.NoError(t, err)
	require.Equal(t, a.Title, "Test title 1")

	// invalid uuid
	_, err = r.GetArticleById("invalid uuid")
	require.Error(t, err)

	// non-existing article
	_, err = r.GetArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
	require.ErrorIs(t, err, ErrArticleNotFound)
}

func TestGetAuthorById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	dumpTestData(t, db)

	// existing author
	a, err := r.GetAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403124")
	require.NotEmpty(t, a)
	require.NoError(t, err)
	require.Equal(t, a.Name, "Test Author1")

	// invalid uuid
	_, err = r.GetAuthorById("invalid uuid")
	require.Error(t, err)

	// non-existing author
	_, err = r.GetAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
	require.ErrorIs(t, err, ErrAuthorNotFound)
}

func TestGetAuthorsByIds(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	dumpTestData(t, db)

	// existing authors
	ids := []string{"b4a4de9e-2f52-4cf1-8907-3d828d403124", "b4a4de9e-2f52-4cf1-8907-3d828d403125"}
	a, err := r.GetAuthorsByIds(ids)
	require.NotEmpty(t, a)
	require.NoError(t, err)
	require.Len(t, a, 2)

	// invalid uuid
	ids = append(ids, "invalid uuid")
	_, err = r.GetAuthorsByIds(ids)
	require.Error(t, err)

}

func TestGetAuthorByNameAndEmail(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	dumpTestData(t, db)

	// existing author
	a, err := r.GetAuthorByNameAndEmail("Test Author1", "test.author1@email.com")
	require.NotEmpty(t, a)
	require.NoError(t, err)
	require.Equal(t, a.Id, "b4a4de9e-2f52-4cf1-8907-3d828d403124")

	// non-existing author
	_, err = r.GetAuthorByNameAndEmail("John Doe", "john.doe@mail.com")
	require.ErrorIs(t, err, ErrAuthorNotFound)
}

func TestAddAuthor(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	id, err := r.AddAuthor(repo.Author{Name: "John Doe", Email: "john.doe@mail.com"})
	require.NoError(t, err)
	require.Len(t, id, 36)

	a, err := r.ListAuthors()
	require.NoError(t, err)
	require.NotEmpty(t, a)
}

func TestAddArticle(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}
	dumpTestData(t, db)

	// author id already in the table, returns no error
	id, err := r.AddArticle(repo.Article{Title: "test", Body: "test", Author: repo.Author{Id: "b4a4de9e-2f52-4cf1-8907-3d828d403124"}})
	require.NoError(t, err)
	require.Len(t, id, 36)

	a, err := r.ListArticles()
	require.NoError(t, err)
	require.Len(t, a, 3)

	// author id not in the table, returns error
	_, err = r.AddArticle(repo.Article{Title: "test", Body: "test", Author: repo.Author{Id: "b4a4de9e-2f52-4cf1-8907-3d828d403128"}})
	require.Error(t, err)
}

func TestDeleteArticleById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}
	dumpTestData(t, db)

	// existing article
	err := r.DeleteArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403126")
	require.NoError(t, err)

	// non-existing article
	err = r.DeleteArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
	require.ErrorIs(t, err, ErrArticleNotFound)

	// closed connection
	db.Close()
	err = r.DeleteArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403126")
	require.Error(t, err)
}

func TestDeleteAuthorById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}
	dumpTestData(t, db)

	// existing article
	err := r.DeleteAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403124")
	require.NoError(t, err)

	// non-existing article
	err = r.DeleteAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
	require.ErrorIs(t, err, ErrAuthorNotFound)

	// closed connection
	db.Close()
	err = r.DeleteAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403124")
	require.Error(t, err)
}

func TestDeleteAuthorByNameAndEmail(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}
	dumpTestData(t, db)

	// existing article
	err := r.DeleteAuthorByNameAndEmail("Test Author1", "test.author1@email.com")
	require.NoError(t, err)

	// non-existing article
	err = r.DeleteAuthorByNameAndEmail("Do not exist", "Do not exist")
	require.ErrorIs(t, err, ErrAuthorNotFound)

	// closed connection
	db.Close()
	err = r.DeleteAuthorByNameAndEmail("Test Author1", "test.author1@email.com")
	require.Error(t, err)
}
