package entities

import "time"

type Comment struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `json:"author,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" validate:"required"`
	AuthorID uint   `json:"author_id" validate:"required"`
	Content  string `json:"content" validate:"required"`
}
