package repository

import (
	"github.com/google/uuid"
)

type UserService interface {
	Get(login string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByPhoneNumber(phone_number string) (*User, error)
	UpdateAuthToken(id uuid.UUID, token string) error
	ValidatePassword(username, password string) (*User, error)
	Save(*User) (*User, error)
}
