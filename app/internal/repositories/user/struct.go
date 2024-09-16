package user

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primary_key"`
	Name         string    `gorm:"type:varchar(100)"`
	Email        string    `gorm:"unique;not null;uniqueIndex"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
