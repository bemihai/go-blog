package blog

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	PostedAt time.Time `json:"posted_at"`
	Author   Author    `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
