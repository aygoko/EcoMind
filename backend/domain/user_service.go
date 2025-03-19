package repository

import 

type UserService interface {
	Get(login string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByPhoneNumber(phone_number string) (*User, error)
	UpdateAuthToken(id uuid.UUID, token string) error
	ValidatePassword(username, password string) (*User, error)
}

