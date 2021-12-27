package postgres

import (
	repo "blog/repo"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type PSQLRepository struct {
	DB *sql.DB
}

// Get all articles.
func (r *PSQLRepository) ListArticles() ([]repo.Article, error) {

	articles := make([]repo.Article, 0)
	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []repo.Article{}, fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var art repo.Article
		var auth repo.Author
		err := rows.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Id)
		if err != nil {
			return []repo.Article{}, fmt.Errorf("cannot scan article: %w", err)
		}
		art.Author = auth
		articles = append(articles, art)
	}
	return articles, nil
}

// Get all authors.
func (r *PSQLRepository) ListAuthors() ([]repo.Author, error) {

	authors := make([]repo.Author, 0)
	query := `SELECT a.id, a.name, a.email FROM blog.authors a;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []repo.Author{}, fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a repo.Author
		err := rows.Scan(&a.Id, &a.Name, &a.Email)
		if err != nil {
			return []repo.Author{}, fmt.Errorf("cannot scan author: %w", err)
		}
		authors = append(authors, a)
	}
	return authors, nil
}

var ErrArticleNotFound = errors.New("article not found")

// Get article by id.
func (r *PSQLRepository) GetArticleById(id string) (repo.Article, error) {

	var art repo.Article
	var auth repo.Author

	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a WHERE a.id = $1;`
	row := r.DB.QueryRow(query, id)

	switch err := row.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Id); err {
	case sql.ErrNoRows:
		return repo.Article{}, ErrArticleNotFound
	case nil:
		art.Author = auth
		return art, nil
	default:
		return repo.Article{}, fmt.Errorf("cannot scan article: %w", err)
	}
}

var ErrAuthorNotFound = errors.New("author not found")

// Get author by id.
func (r *PSQLRepository) GetAuthorById(id string) (repo.Author, error) {

	var a repo.Author

	query := `SELECT a.id, a.name, a.email FROM blog.authors a WHERE a.id = $1;`
	row := r.DB.QueryRow(query, id)

	switch err := row.Scan(&a.Id, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return repo.Author{}, ErrAuthorNotFound
	case nil:
		return a, nil
	default:
		return repo.Author{}, fmt.Errorf("cannot scan author: %w", err)
	}
}

// Get authors by ids.
func (r *PSQLRepository) GetAuthorsByIds(ids []string) ([]repo.Author, error) {

	authors := make([]repo.Author, 0)

	query := `SELECT a.id, a.name, a.email FROM blog.authors a WHERE a.id = any($1);`
	rows, err := r.DB.Query(query, pq.Array(ids))
	if err != nil {
		return []repo.Author{}, fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a repo.Author
		err := rows.Scan(&a.Id, &a.Name, &a.Email)
		if err != nil {
			return []repo.Author{}, fmt.Errorf("cannot scan author: %w", err)
		}
		authors = append(authors, a)
	}

	return authors, nil
}

// Get author by name and email.
func (r *PSQLRepository) GetAuthorByNameAndEmail(name string, email string) (repo.Author, error) {

	var a repo.Author

	query := `SELECT a.id, a.name, a.email FROM blog.authors a WHERE a.name = $1 AND a.email = $2;`
	row := r.DB.QueryRow(query, name, email)

	switch err := row.Scan(&a.Id, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return repo.Author{}, ErrAuthorNotFound
	case nil:
		return a, nil
	default:
		return repo.Author{}, fmt.Errorf("cannot scan author: %w", err)
	}
}

// Add new author and return its id.
func (r *PSQLRepository) AddAuthor(a repo.Author) (string, error) {

	var id string

	// TO DO: email should be unique, return error if exists
	query := `INSERT INTO blog.authors(name, email) values ($1, $2) RETURNING id;`
	err := r.DB.QueryRow(query, a.Name, a.Email).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("cannot execute query: %w", err)
	}

	return id, nil
}

// Add new article and return its id.
func (r *PSQLRepository) AddArticle(a repo.Article) (string, error) {

	var id string

	query := `INSERT INTO blog.articles(title, body, author_id) values ($1, $2, $3) RETURNING id;`
	err := r.DB.QueryRow(query, a.Title, a.Body, a.Author.Id).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("cannot execute query: %w", err)
	}

	return id, nil
}

// Delete article by id.
func (r *PSQLRepository) DeleteArticleById(id string) error {

	query := `DELETE FROM blog.articles WHERE id = $1;`
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
func (r *PSQLRepository) DeleteAuthorById(id string) error {

	query := `DELETE FROM blog.authors WHERE id = $1;`
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
func (r *PSQLRepository) DeleteAuthorByNameAndEmail(name string, email string) error {

	query := `DELETE FROM blog.authors WHERE name = $1 AND email = $2;`
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
