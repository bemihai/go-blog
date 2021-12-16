package blog

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var auth = &Author{
	Name:  "John Doe",
	Email: "john.doe@email.com",
}

var art = &Article{
	Id:       uuid.New(),
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

	listArticlesQuery := `SELECT art.id, art.title, art.body, art.posted_at, auth.name, auth.email
						FROM blog.articles art
						LEFT JOIN blog.authors auth on art.author_id = auth.id;`

	rows := sqlmock.NewRows([]string{"id", "title", "body", "posted_at", "name", "email"}).
		AddRow(art.Id, art.Title, art.Body, art.PostedAt, art.Author.Name, art.Author.Email).
		AddRow(art.Id, art.Title, art.Body, art.PostedAt, art.Author.Name, art.Author.Email)

	mock.ExpectQuery(listArticlesQuery).WillReturnRows(rows)

	articles, err := repo.ListArticles()

	assert.NotNil(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
}

func TestListArticlesEmpty(t *testing.T) {

	db, mock := NewMock()
	repo := &PSQLRepository{db}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "body", "posted_at", "name", "email"})

	mock.ExpectQuery(listArticlesQuery).WillReturnRows(rows)

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
	mock.ExpectQuery(listArticlesQuery).WillReturnError(dbError)

	_, err := repo.ListArticles()
	targetErr := &DatabaseError{}

	assert.Error(t, err)
	assert.ErrorAs(t, err, targetErr)

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
