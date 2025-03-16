package repository

import (
	"errors"

	repository "github.com/aygoko/EcoMInd/backend/domain"
)

func NewUserRepository() repository.UserService {
	return &UserRepositoryRAM{
		data: make(map[string]*repository.User),
	}
}

type UserRepositoryRAM struct {
	data map[string]*repository.User
}

func (r *UserRepositoryRAM) Get(login string) (*repository.User, error) {
	user, exists := r.data[login]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepositoryRAM) GetByEmail(email string) (*repository.User, error) {
	for _, user := range r.data {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepositoryRAM) GetByPhoneNumber(phone_number string) (*repository.User, error) {
	for _, user := range r.data {
		if user.PhoneNumber == phone_number {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
