package repo

import (
	"time"
)

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
