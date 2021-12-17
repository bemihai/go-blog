package repository

import (
	"time"

	"github.com/google/uuid"
)

// Repository represents the blog repository.
type Repository interface {
	ListArticles() ([]Article, error)
	ListAuthors() ([]Author, error)
	GetArticleById(id uuid.UUID) (Article, error)
	GetAuthorById(id uuid.UUID) (Author, error)
	GetAuthorByNameAndEmail(name string, email string) (Author, error)
	AddArticle(a Article) (uuid.UUID, error)
	AddAuthor(a Author) (uuid.UUID, error)
	DeleteArticleById(id uuid.UUID) error
	DeleteAuthorById(id uuid.UUID) error
	DeleteAuthorByNameAndEmail(name string, email string) error
}

// Article represents the article model.
type Article struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	PostedAt time.Time `json:"posted_at"`
	Author   Author    `json:"author"`
}

// Author represents the author model.
type Author struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
