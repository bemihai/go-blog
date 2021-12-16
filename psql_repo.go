package blog

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PSQLRepository struct {
	DB *sql.DB
}

// List all articles from the DB.
func (repo *PSQLRepository) ListArticles() ([]Article, error) {

	articles := make([]Article, 0)

	rows, err := repo.DB.Query(listArticlesQuery)

	if err != nil {
		return articles, NewDatabaseError(err, "cannot fetch articles")
	}
	defer rows.Close()

	for rows.Next() {

		var art Article
		var auth Author

		err := rows.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Name, &auth.Email)
		if err != nil {
			return articles, NewDatabaseError(err, "cannot scan article")
		}

		art.Author = auth
		articles = append(articles, art)
	}

	return articles, nil
}

// Get article by id.
func (repo *PSQLRepository) GetArticleById(id uuid.UUID) (Article, error) {

	var art Article
	var auth Author

	row := repo.DB.QueryRow(getArticleByIdQuery, id)

	switch err := row.Scan(&art.Id, &art.Title, &art.Body, &art.PostedAt, &auth.Name, &auth.Email); err {
	case sql.ErrNoRows:
		return art, NewDatabaseError(err, fmt.Sprintf("article with id %s does not exist", id))
	case nil:
		art.Author = auth
		return art, nil
	default:
		return art, NewDatabaseError(err, "database error")
	}
}

// Add new author to the DB and return its id.
func (repo *PSQLRepository) addNewAuthor(a Author) (uuid.UUID, error) {

	var id uuid.UUID

	checkAuthorQuery := `SELECT a.id FROM blog.authors a WHERE a.name = $1 AND a.email = $2;`
	addAuthorQuery := `INSERT INTO blog.authors(name, email) values ($1, $2) RETURNING id;`

	// check if author already exists in the DB by name and email
	row := repo.DB.QueryRow(checkAuthorQuery, a.Name, a.Email)

	switch err_1 := row.Scan(&id); err_1 {

	// if author not found, add it and return the id
	case sql.ErrNoRows:
		err_2 := repo.DB.QueryRow(addAuthorQuery, a.Name, a.Email).Scan(&id)
		if err_2 != nil {
			return id, fmt.Errorf("cannot add new author, %w", err_2)
		}
		return id, nil
	// if author found, return id
	case nil:
		return id, nil
	default:
		return id, fmt.Errorf("database error, %w", err_1)
	}

}

// Add new article to the DB.
func (repo *PSQLRepository) PostArticle(a Article) (Article, error) {

	var articleId uuid.UUID

	// add author and retrieve its id
	authorId, err := repo.addNewAuthor(a.Author)
	if err != nil {
		return Article{}, err
	}

	// add article and retrieve its id
	err = repo.DB.QueryRow(addArticleQuery, a.Title, a.Body, authorId).Scan(&articleId)
	if err != nil {
		return Article{}, fmt.Errorf("database error, %w", err)
	}

	// retrieve article from DB and return it
	var article Article
	article, err_1 := repo.GetArticleById(articleId)
	if err_1 != nil {
		return Article{}, fmt.Errorf("database error, %w", err_1)
	}

	return article, nil
}

// Delete article by id.
func (repo *PSQLRepository) DeleteArticleById(id uuid.UUID) (Article, error) {

	var a Article
	var authorId uuid.UUID

	deleteArticleQuery := `DELETE FROM blog.articles WHERE id = $1 RETURNING *;`

	row := repo.DB.QueryRow(deleteArticleQuery, id)

	switch err := row.Scan(&a.Id, &a.Title, &a.Body, &a.PostedAt, &authorId); err {
	case sql.ErrNoRows:
		return a, NewDatabaseError(err, fmt.Sprintf("article with id %s does not exist", id))
	case nil:
		return a, nil
	default:
		return a, NewDatabaseError(err, "database error")
	}
}

// Delete author by name and email (and all its articles).
func (repo *PSQLRepository) DeleteAuthorByNameAndEmail(name string, email string) (Author, error) {

	var a Author
	var authorId uuid.UUID

	deleteAuthorQuery := `DELETE FROM blog.authors WHERE name = $1 AND email = $2 RETURNING *;`

	row := repo.DB.QueryRow(deleteAuthorQuery, name, email)
	switch err := row.Scan(&authorId, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		return a, fmt.Errorf("Author with name %s does not exist, %w", name, err)
	case nil:
		return a, nil
	default:
		return a, fmt.Errorf("database error, %w", err)
	}
}

// Database errors wrapper.
type DatabaseError struct {
	Message string
	Err     error
}

// Implement the error interface.
func (err DatabaseError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

// Database error constructor.
func NewDatabaseError(err error, message string) error {
	return DatabaseError{
		Message: message,
		Err:     err,
	}
}
