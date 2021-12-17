package postgres

import (
	repo "blog/repo"

	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var auth = &repo.Author{
	Name:  "John Doe",
	Email: "john.doe@email.com",
}

var art = &repo.Article{
	Id:       "uuid.New()",
	Title:    "Test title",
	Body:     "Test body",
	PostedAt: time.Now(),
	Author:   *auth,
}

// Mock database connection.
func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestListArticles(t *testing.T) {

	db, mock := NewMock()
	repo := &PSQLRepository{db}
	defer db.Close()

	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a;`

	rows := sqlmock.NewRows([]string{"id", "title", "body", "posted_at", "name", "email"}).
		AddRow(art.Id, art.Title, art.Body, art.PostedAt, art.Author.Name, art.Author.Email).
		AddRow(art.Id, art.Title, art.Body, art.PostedAt, art.Author.Name, art.Author.Email)

	mock.ExpectQuery(query).WillReturnRows(rows)

	articles, err := repo.ListArticles()

	assert.NotNil(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
}

func TestListArticlesEmpty(t *testing.T) {

	db, mock := NewMock()
	repo := &PSQLRepository{db}
	defer db.Close()

	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a;`
	rows := sqlmock.NewRows([]string{"id", "title", "body", "posted_at", "name", "email"})

	mock.ExpectQuery(query).WillReturnRows(rows)

	articles, err := repo.ListArticles()

	assert.Empty(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 0)

}

func TestListArticlesDatabaseError(t *testing.T) {

	db, mock := NewMock()
	repo := &PSQLRepository{db}
	defer db.Close()

	dbError := errors.New("database error")
	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a;`
	mock.ExpectQuery(query).WillReturnError(dbError)

	_, err := repo.ListArticles()
	assert.Error(t, err)

}

// TODO: simplify queries in repository, it's hard to test multiple queries in one call
// func TestPostArticle(t *testing.T) {

// 	db, mock := NewMock()
// 	repo := &PSQLRepository{db}
// 	defer db.Close()

// 	prep := mock.ExpectPrepare(addArticleQuery)
// 	prep.ExpectExec().WithArgs(art).WillReturnResult(sqlmock.NewResult(1, 0))

// 	_, err := repo.PostArticle(*art)

// 	assert.NoError(t, err)
// }
