package postgres

import (
	repo "blog/repo"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PSQLRepository struct {
	DB *sql.DB
}

// List all articles.
func (r *PSQLRepository) ListArticles() ([]repo.Article, error) {

	articles := make([]repo.Article, 0)
	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a;`

	rows, err := r.DB.Query(query)

	if err != nil {
		return articles, NewDatabaseError(err, "cannot fetch articles")
	}

	defer rows.Close()

	for rows.Next() {

		var art repo.Article
		var auth repo.Author

		err := rows.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Id)

		if err != nil {
			return articles, NewDatabaseError(err, "cannot scan article")
		}

		art.Author = auth
		articles = append(articles, art)
	}

	return articles, nil
}

// List all authors.
func (r *PSQLRepository) ListAuthors() ([]repo.Author, error) {

	authors := make([]repo.Author, 0)
	query := `SELECT a.id, a.name, a.email FROM blog.authors a;`

	rows, err := r.DB.Query(query)

	if err != nil {
		return authors, NewDatabaseError(err, "cannot fetch authors")
	}

	defer rows.Close()

	for rows.Next() {

		var a repo.Author

		err := rows.Scan(&a.Id, &a.Name, &a.Email)

		if err != nil {
			return authors, NewDatabaseError(err, "cannot scan article")
		}

		authors = append(authors, a)
	}

	return authors, nil
}

// Get article by id.
func (r *PSQLRepository) GetArticleById(id uuid.UUID) (repo.Article, error) {

	var art repo.Article
	var auth repo.Author
	query := `SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM blog.articles a WHERE a.id = $1;`

	row := r.DB.QueryRow(query, id)

	switch err := row.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Id); err {
	case sql.ErrNoRows:
		return art, NewDatabaseError(err, fmt.Sprintf("article with id '%s' does not exist", id))
	case nil:
		art.Author = auth
		return art, nil
	default:
		return art, NewDatabaseError(err, "database error")
	}
}

// Get author by id.
func (r *PSQLRepository) GetAuthorById(id uuid.UUID) (repo.Author, error) {

	var a repo.Author
	query := `SELECT a.id, a.name, a.email FROM blog.authors a WHERE a.id = $1;`

	row := r.DB.QueryRow(query, id)

	switch err := row.Scan(&a.Id, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return a, NewDatabaseError(err, fmt.Sprintf("article with id '%s' does not exist", id))
	case nil:
		return a, nil
	default:
		return a, NewDatabaseError(err, "database error")
	}
}

// Get author by name and email.
func (r *PSQLRepository) GetAuthorByNameAndEmail(name string, email string) (repo.Author, error) {

	var a repo.Author
	query := `SELECT a.id FROM blog.authors a WHERE a.name = $1 AND a.email = $2;`

	row := r.DB.QueryRow(query, name, email)

	switch err := row.Scan(&a.Id, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return a, NewDatabaseError(err, fmt.Sprintf("article with name '%s' and email '%s' does not exist", name, email))
	case nil:
		return a, nil
	default:
		return a, NewDatabaseError(err, "database error")
	}

}

// Add new author and return its id.
func (r *PSQLRepository) AddAuthor(a repo.Author) (uuid.UUID, error) {

	var id uuid.UUID
	query := `INSERT INTO blog.authors(name, email) values ($1, $2) RETURNING id;`

	err := r.DB.QueryRow(query, a.Name, a.Email).Scan(&id)

	if err != nil {
		return id, NewDatabaseError(err, "cannot add new author")
	}

	return id, nil
}

// Add new article and return its id.
func (r *PSQLRepository) AddArticle(a repo.Article) (uuid.UUID, error) {

	var id uuid.UUID
	query := `INSERT INTO blog.articles(title, body, author_id) values ($1, $2, $3) RETURNING id;`

	// add article and retrieve its id
	err := r.DB.QueryRow(query, a.Title, a.Body, a.Author.Id).Scan(&id)
	if err != nil {
		return id, NewDatabaseError(err, "cannot add new article")
	}

	return id, nil
}

// Delete article by id.
func (r *PSQLRepository) DeleteArticleById(id uuid.UUID) error {

	query := `DELETE FROM blog.articles WHERE id = $1;`

	res, err := r.DB.Exec(query, id)
	if err != nil {
		return NewDatabaseError(err, "database error")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return NewDatabaseError(err, "database error")
	}

	if count == 0 {
		return NewDatabaseError(err, fmt.Sprintf("article with id '%s' does not exist", id))
	}
	return nil
}

// Delete author by id.
func (r *PSQLRepository) DeleteAuthorById(id uuid.UUID) error {

	query := `DELETE FROM blog.authors WHERE id = $1;`

	res, err := r.DB.Exec(query, id)
	if err != nil {
		return NewDatabaseError(err, "database error")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return NewDatabaseError(err, "database error")
	}

	if count == 0 {
		return NewDatabaseError(err, fmt.Sprintf("author with id '%s' does not exist", id))
	}
	return nil
}

// Delete author by name and email (and all its articles).
func (r *PSQLRepository) DeleteAuthorByNameAndEmail(name string, email string) error {

	query := `DELETE FROM blog.authors WHERE name = $1 AND email = $2;`

	res, err := r.DB.Exec(query, name, email)
	if err != nil {
		return NewDatabaseError(err, "database error")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return NewDatabaseError(err, "database error")
	}

	if count == 0 {
		return NewDatabaseError(err, fmt.Sprintf("author with name '%s' and email '%s' does not exist", name, email))
	}
	return nil
}

// DatabaseError implements the error interface and wrapps database errors.
type DatabaseError struct {
	Message string
	Err     error
}

func (err DatabaseError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

// DatabaseError constructor.
func NewDatabaseError(err error, message string) error {
	return DatabaseError{
		Message: message,
		Err:     err,
	}
}
