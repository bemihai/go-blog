// Code generated by sqlc. DO NOT EDIT.
// source: articles.sql

package sqlc_db

import (
	"context"

	"github.com/google/uuid"
)

const deleteArticleById = `-- name: DeleteArticleById :exec
DELETE FROM articles WHERE id = $1
`

func (q *Queries) DeleteArticleById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteArticleById, id)
	return err
}

const getArticleById = `-- name: GetArticleById :one
SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM articles a WHERE a.id = $1
`

func (q *Queries) GetArticleById(ctx context.Context, id uuid.UUID) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticleById, id)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.PostedAt,
		&i.AuthorID,
	)
	return i, err
}

const listArticles = `-- name: ListArticles :many
SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM articles a
`

func (q *Queries) ListArticles(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.PostedAt,
			&i.AuthorID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAuthors = `-- name: ListAuthors :many
SELECT a.id, a.name, a.email FROM authors a
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.QueryContext(ctx, listAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAuthor = `-- name: UpdateAuthor :exec
UPDATE authors SET name = $2 AND email = $3 WHERE id = $1
`

type UpdateAuthorParams struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) error {
	_, err := q.db.ExecContext(ctx, updateAuthor, arg.ID, arg.Name, arg.Email)
	return err
}