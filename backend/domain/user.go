package repository

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Username     string    `json:"login" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PhoneNumber  string    `json:"phone_number"`
	PasswordHash string    `json:"-" gorm:"not null"`
	AuthToken    string    `json:"auth_token"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
