package post

import (
	"time"
)

type Post struct {
	ID        uint      `gorm:"primary_key"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	AuthorID  uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
