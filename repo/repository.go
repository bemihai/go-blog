package repository

import (
	"time"
)

// Repository represents the blog repository.
type Repository interface {
	ListArticles() ([]Article, error)
	ListAuthors() ([]Author, error)
	GetArticleById(id string) (Article, error)
	GetAuthorById(id string) (Author, error)
	GetAuthorsByIds(ids []string) ([]Author, error)
	GetAuthorByNameAndEmail(name string, email string) (Author, error)
	AddArticle(a Article) (string, error)
	AddAuthor(a Author) (string, error)
	DeleteArticleById(id string) error
	DeleteAuthorById(id string) error
	DeleteAuthorByNameAndEmail(name string, email string) error
}

// Article represents the article model.
type Article struct {
	Id       string    `json:"-"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	PostedAt time.Time `json:"posted_at"`
	Author   Author    `json:"author"`
}

// Author represents the author model.
type Author struct {
	Id    string `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
