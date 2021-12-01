package blog

import (
	"database/sql"
	"fmt"
)

type PSQLRepository struct {
	DB *sql.DB
}

// List all articles from the DB.
func (repo *PSQLRepository) ListArticles() ([]Article, error) {

	articles := make([]Article, 0)
	rows, err := repo.DB.Query(`SELECT * FROM blog.articles;`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var a Article
		var authorId int
		err := rows.Scan(&a.Id, &a.Title, &a.Body, &a.PostedAt, &authorId)
		if err != nil {
			panic(err)
		}
		// for each article, get author by id from blog.authors table
		author, err := repo.getAuthorById(authorId)
		if err != nil {
			panic(err)
		}
		a.Author = author
		articles = append(articles, a)
	}
	return articles, nil
}

// Get author by id.
func (repo *PSQLRepository) getAuthorById(id int) (Author, error) {

	getAuthorByIdQuery := `SELECT a.name, a.email from blog.authors a WHERE a.id = $1;`
	var author Author
	row := repo.DB.QueryRow(getAuthorByIdQuery, id)
	switch err := row.Scan(&author.Name, &author.Email); err {
	case sql.ErrNoRows:
		err := fmt.Errorf("Author with id %d does not exist", id)
		panic(err)
	case nil:
		return author, nil
	default:
		panic(err)
	}
}

// Get article by id.
func (repo *PSQLRepository) GetArticleById(id int) (Article, error) {

	getArticleByIdQuery := `SELECT * from blog.articles a WHERE a.id = $1;`
	var a Article
	var authorId int
	row := repo.DB.QueryRow(getArticleByIdQuery, id)
	switch err := row.Scan(&a.Id, &a.Title, &a.Body, &a.PostedAt, &authorId); err {
	case sql.ErrNoRows:
		err := fmt.Errorf("Article with id %d does not exist", id)
		panic(err)
	case nil:
		// get author by id from blog.authors table
		author, err := repo.getAuthorById(authorId)
		if err != nil {
			panic(err)
		}
		a.Author = author
		return a, nil
	default:
		panic(err)
	}
}

// Add new author to the DB and return its id.
func (repo *PSQLRepository) addNewAuthor(a Author) (int, error) {

	checkAuthorQuery := `SELECT * FROM blog.authors a WHERE a.name = $1 AND a.email = $2;`
	addAuthorQuery := `INSERT INTO blog.authors(name, email) values ($1, $2) RETURNING id;`
	var id int
	var author Author

	// check if author already exists in the DB by name and email
	row := repo.DB.QueryRow(checkAuthorQuery, a.Name, a.Email)

	switch err := row.Scan(&id, &author.Name, &author.Email); err {
	// if author not found, add it and return the id
	case sql.ErrNoRows:
		err_ := repo.DB.QueryRow(addAuthorQuery, a.Name, a.Email).Scan(&id)
		if err_ != nil {
			panic(err)
		}
		return id, nil
	// if author found, return id
	case nil:
		return id, nil
	default:
		panic(err)
	}

}

// Add new article to the DB.
func (repo *PSQLRepository) PostArticle(a Article) (Article, error) {

	addArticleQuery := `INSERT INTO blog.articles(title, body, author_id) values ($1, $2, $3) RETURNING id;`
	var articleId int

	// add author and retrieve its id
	authorId, err := repo.addNewAuthor(a.Author)
	if err != nil {
		panic(err)
	}

	// add article and retrieve its id
	err = repo.DB.QueryRow(addArticleQuery, a.Title, a.Body, authorId).Scan(&articleId)
	if err != nil {
		panic(err)
	}

	// retrieve article from DB and return it
	var article Article
	article, err_ := repo.GetArticleById(articleId)
	if err_ != nil {
		panic(err)
	}

	return article, nil
}

// Delete article by id.
func (repo *PSQLRepository) DeleteArticleById(id int) (Article, error) {

	deleteArticleQuery := `DELETE FROM blog.articles WHERE id = $1 RETURNING *;`
	var a Article
	var authorId int

	row := repo.DB.QueryRow(deleteArticleQuery, id)
	switch err := row.Scan(&a.Id, &a.Title, &a.Body, &a.PostedAt, &authorId); err {
	case sql.ErrNoRows:
		err := fmt.Errorf("Article with id %d does not exist", id)
		panic(err)
	case nil:
		author, err := repo.getAuthorById(authorId)
		if err != nil {
			panic(err)
		}
		a.Author = author
		return a, nil
	default:
		panic(err)
	}
}

// Delete author by name and email (and all its articles).
func (repo *PSQLRepository) DeleteAuthorByNameAndEmail(name string, email string) (Author, error) {

	deleteAuthorQuery := `DELETE FROM blog.authors WHERE name = $1 AND email = $2 RETURNING *;`
	var a Author
	var authorId int

	row := repo.DB.QueryRow(deleteAuthorQuery, name, email)
	switch err := row.Scan(&authorId, &a.Name, &a.Email); err {
	case sql.ErrNoRows:
		err := fmt.Errorf("Author with name %s does not exist", name)
		panic(err)
	case nil:
		return a, nil
	default:
		panic(err)
	}
}
