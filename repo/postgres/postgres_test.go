package postgres

import (
	repo "blog/repo"

	"testing"

	"github.com/stretchr/testify/require"
)

const (
	connection = "postgres://localhost:5432/blog?user=postgres&password=postgres&sslmode=disable"
)

func TestListArticles(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("table containing 2 entries", func(t *testing.T) {
		articles, err := r.ListArticles()
		require.NotEmpty(t, articles)
		require.NoError(t, err)
		require.Len(t, articles, 2)
	})

	t.Run("empty table", func(t *testing.T) {
		truncateTables(t, db)
		articles, err := r.ListArticles()
		require.Empty(t, articles)
		require.NoError(t, err)
		require.Len(t, articles, 0)
	})
}

func TestListAuthors(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("table containing 2 entries", func(t *testing.T) {
		authors, err := r.ListAuthors()
		require.NotEmpty(t, authors)
		require.NoError(t, err)
		require.Len(t, authors, 2)
	})

	t.Run("empty table", func(t *testing.T) {
		truncateTables(t, db)
		authors, err := r.ListAuthors()
		require.Empty(t, authors)
		require.NoError(t, err)
		require.Len(t, authors, 0)
	})
}

func TestGetArticleById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing article", func(t *testing.T) {
		a, err := r.GetArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403126")
		require.NotEmpty(t, a)
		require.NoError(t, err)
		require.Equal(t, a.Title, "Test title 1")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := r.GetArticleById("invalid uuid")
		require.Error(t, err)
	})

	t.Run("non-existing article", func(t *testing.T) {
		_, err := r.GetArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
		require.ErrorIs(t, err, ErrArticleNotFound)
	})
}

func TestGetAuthorById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing author", func(t *testing.T) {
		a, err := r.GetAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403124")
		require.NotEmpty(t, a)
		require.NoError(t, err)
		require.Equal(t, a.Name, "Test Author1")
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := r.GetAuthorById("invalid uuid")
		require.Error(t, err)
	})

	t.Run("non-existing author", func(t *testing.T) {
		_, err := r.GetAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
		require.ErrorIs(t, err, ErrAuthorNotFound)
	})
}

func TestGetAuthorsByIds(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing authors", func(t *testing.T) {
		ids := []string{"b4a4de9e-2f52-4cf1-8907-3d828d403124", "b4a4de9e-2f52-4cf1-8907-3d828d403125"}
		a, err := r.GetAuthorsByIds(ids)
		require.NotEmpty(t, a)
		require.NoError(t, err)
		require.Len(t, a, 2)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		ids := []string{"b4a4de9e-2f52-4cf1-8907-3d828d403124", "invalid uuid"}
		_, err := r.GetAuthorsByIds(ids)
		require.Error(t, err)
	})

}

func TestGetAuthorByNameAndEmail(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing author", func(t *testing.T) {
		a, err := r.GetAuthorByNameAndEmail("Test Author1", "test.author1@email.com")
		require.NotEmpty(t, a)
		require.NoError(t, err)
		require.Equal(t, a.Id, "b4a4de9e-2f52-4cf1-8907-3d828d403124")
	})

	t.Run("non-existing author", func(t *testing.T) {
		_, err := r.GetAuthorByNameAndEmail("John Doe", "john.doe@mail.com")
		require.ErrorIs(t, err, ErrAuthorNotFound)
	})
}

func TestAddAuthor(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("valid author", func(t *testing.T) {
		id, err := r.AddAuthor(repo.Author{Name: "John Doe", Email: "john.doe@mail.com"})
		require.NoError(t, err)
		require.Len(t, id, 36)
		a, err := r.ListAuthors()
		require.NoError(t, err)
		require.NotEmpty(t, a)
	})
}

func TestAddArticle(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("author id already in the table", func(t *testing.T) {
		id, err := r.AddArticle(repo.Article{Title: "test", Body: "test", Author: repo.Author{Id: "b4a4de9e-2f52-4cf1-8907-3d828d403124"}})
		require.NoError(t, err)
		require.Len(t, id, 36)
		a, err := r.ListArticles()
		require.NoError(t, err)
		require.Len(t, a, 3)
	})

	t.Run("author id not in the table", func(t *testing.T) {
		_, err := r.AddArticle(repo.Article{Title: "test", Body: "test", Author: repo.Author{Id: "b4a4de9e-2f52-4cf1-8907-3d828d403128"}})
		require.Error(t, err)
	})
}

func TestDeleteArticleById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing article", func(t *testing.T) {
		err := r.DeleteArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403126")
		require.NoError(t, err)
	})

	t.Run("non-existing article", func(t *testing.T) {
		err := r.DeleteArticleById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
		require.ErrorIs(t, err, ErrArticleNotFound)
	})
}

func TestDeleteAuthorById(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing author", func(t *testing.T) {
		err := r.DeleteAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403124")
		require.NoError(t, err)
	})

	t.Run("non-existing author", func(t *testing.T) {
		err := r.DeleteAuthorById("b4a4de9e-2f52-4cf1-8907-3d828d403128")
		require.ErrorIs(t, err, ErrAuthorNotFound)
	})
}

func TestDeleteAuthorByNameAndEmail(t *testing.T) {

	db, _ := createTestDB(t, connection)
	r := PSQLRepository{DB: db}

	t.Run("existing article", func(t *testing.T) {
		err := r.DeleteAuthorByNameAndEmail("Test Author1", "test.author1@email.com")
		require.NoError(t, err)
	})

	t.Run("non-existing article", func(t *testing.T) {
		err := r.DeleteAuthorByNameAndEmail("Do not exist", "Do not exist")
		require.ErrorIs(t, err, ErrAuthorNotFound)
	})
}
