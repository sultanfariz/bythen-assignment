package entities

import (
	"time"
)

type Post struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	AuthorID uint   `json:"author_id" validate:"required"`
}

type UpdatePostRequest struct {
	ID      uint   `json:"id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
