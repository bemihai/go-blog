-- name: ListArticles :many
SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM articles a;

-- name: ListAuthors :many
SELECT a.id, a.name, a.email FROM authors a;

-- name: GetArticleById :one
SELECT a.id, a.title, a.body, a.posted_at, a.author_id FROM articles a WHERE a.id = $1;

-- name: UpdateAuthor :exec
UPDATE authors SET name = $2 AND email = $3 WHERE id = $1;

-- name: DeleteArticleById :exec
DELETE FROM articles WHERE id = $1;