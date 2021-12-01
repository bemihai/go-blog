package blog

import "time"

type Article struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	PostedAt time.Time `json:"posted_at"`
	Author   Author    `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
