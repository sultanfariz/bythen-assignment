package user

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primary_key"`
	Name         string `gorm:"type:varchar(100)"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
