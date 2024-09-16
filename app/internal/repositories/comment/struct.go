package comment

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primary_key"`
	PostID    uint      `gorm:"not null;index"`
	AuthorID  uint      `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}
