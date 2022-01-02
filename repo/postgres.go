package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// BlogService represents the blog repository.
type BlogService interface {
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

type BlogRepo struct {
	DB *sql.DB
}

// Get all articles.
func (r *BlogRepo) ListArticles() ([]Article, error) {

	articles := make([]Article, 0)
	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM articles a;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []Article{}, fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var art Article
		var auth Author
		err := rows.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Id)
		if err != nil {
			return []Article{}, fmt.Errorf("cannot scan article: %w", err)
		}
		art.Author = auth
		articles = append(articles, art)
	}
	return articles, nil
}

// Get all authors.
func (r *BlogRepo) ListAuthors() ([]Author, error) {

	authors := make([]Author, 0)
	query := `SELECT a.id, a.name, a.email FROM authors a;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []Author{}, fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a Author
		err := rows.Scan(&a.Id, &a.Name, &a.Email)
		if err != nil {
			return []Author{}, fmt.Errorf("cannot scan author: %w", err)
		}
		authors = append(authors, a)
	}
	return authors, nil
}

var ErrArticleNotFound = errors.New("article not found")

// Get article by id.
func (r *BlogRepo) GetArticleById(id string) (Article, error) {

	var art Article
	var auth Author

	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM articles a WHERE a.id = $1;`
	row := r.DB.QueryRow(query, id)

	switch err := row.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Id); err {
	case sql.ErrNoRows:
		return Article{}, ErrArticleNotFound
	case nil:
		art.Author = auth
		return art, nil
	default:
		return Article{}, fmt.Errorf("cannot scan article: %w", err)
	}
}

var ErrAuthorNotFound = errors.New("author not found")

// Get author by id.
func (r *BlogRepo) GetAuthorById(id string) (Author, error) {

	var a Author

	query := `SELECT a.id, a.name, a.email FROM authors a WHERE a.id = $1;`
	row := r.DB.QueryRow(query, id)

	switch err := row.Scan(&a.Id, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return Author{}, ErrAuthorNotFound
	case nil:
		return a, nil
	default:
		return Author{}, fmt.Errorf("cannot scan author: %w", err)
	}
}

// Get authors by ids.
func (r *BlogRepo) GetAuthorsByIds(ids []string) ([]Author, error) {

	authors := make([]Author, 0)

	query := `SELECT a.id, a.name, a.email FROM authors a WHERE a.id = any($1);`
	rows, err := r.DB.Query(query, pq.Array(ids))
	if err != nil {
		return []Author{}, fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a Author
		err := rows.Scan(&a.Id, &a.Name, &a.Email)
		if err != nil {
			return []Author{}, fmt.Errorf("cannot scan author: %w", err)
		}
		authors = append(authors, a)
	}

	return authors, nil
}

// Get author by name and email.
func (r *BlogRepo) GetAuthorByNameAndEmail(name string, email string) (Author, error) {

	var a Author

	query := `SELECT a.id, a.name, a.email FROM authors a WHERE a.name = $1 AND a.email = $2;`
	row := r.DB.QueryRow(query, name, email)

	switch err := row.Scan(&a.Id, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return Author{}, ErrAuthorNotFound
	case nil:
		return a, nil
	default:
		return Author{}, fmt.Errorf("cannot scan author: %w", err)
	}
}

// Add new author and return its id.
func (r *BlogRepo) AddAuthor(a Author) (string, error) {

	var id string

	// TO DO: email should be unique, return error if exists
	query := `INSERT INTO authors(name, email) values ($1, $2) RETURNING id;`
	err := r.DB.QueryRow(query, a.Name, a.Email).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("cannot execute query: %w", err)
	}

	return id, nil
}

// Add new article and return its id.
func (r *BlogRepo) AddArticle(a Article) (string, error) {

	var id string

	// author id must exist in the authors table
	query := `INSERT INTO articles(title, body, author_id) values ($1, $2, $3) RETURNING id;`
	err := r.DB.QueryRow(query, a.Title, a.Body, a.Author.Id).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("cannot execute query: %w", err)
	}

	return id, nil
}

// Delete article by id.
func (r *BlogRepo) DeleteArticleById(id string) error {

	query := `DELETE FROM articles WHERE id = $1;`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot retrieve rows affected: %w", err)
	}
	if count == 0 {
		return ErrArticleNotFound
	}

	return nil
}

// Delete author by id.
func (r *BlogRepo) DeleteAuthorById(id string) error {

	query := `DELETE FROM authors WHERE id = $1;`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot retrieve rows affected: %w", err)
	}
	if count == 0 {
		return ErrAuthorNotFound
	}

	return nil
}

// Delete author by name and email (and all its articles).
func (r *BlogRepo) DeleteAuthorByNameAndEmail(name string, email string) error {

	query := `DELETE FROM authors WHERE name = $1 AND email = $2;`
	res, err := r.DB.Exec(query, name, email)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot retrieve rows affected: %w", err)
	}
	if count == 0 {
		return ErrAuthorNotFound
	}

	return nil
}
