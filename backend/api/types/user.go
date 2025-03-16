package types

import (
	"errors"

	"github.com/aygoko/EcoMInd/backend/repository"
	"github.com/aygoko/EcoMInd/backend/usecases/service"
	"golang.org/x/crypto/bcrypt"
)

func (s *repository.UserService) CreateUser(login, email, phone_number, password string) (*repository.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if login == "" || email == "" || phone_number == "" {
		return nil, errors.New("login and email and phone_number are required")
	}

	if _, exists := s.Users[login]; exists {
		return nil, errors.New("user with this login already exists")
	}

	if existingLogin, exists := s.Emails[email]; exists {
		return nil, errors.New("email already in use by " + existingLogin)
	}

	if existingLogin, exists := s.PhoneNumbers[phone_number]; exists {
		return nil, errors.New("phone_number already in use by " + existingLogin)
	}

	user := &repository.User{
		ID:          service.GenerateUserID(),
		Login:       login,
		Email:       email,
		PhoneNumber: phone_number,
		Password:    string(hashedPassword),
	}

	s.Users[login] = user
	s.Emails[email] = login

	return user, nil
}

func (s *repository.UserService) GetUserByLogin(login string) (*repository.User, error) {
	user, exists := s.Users[login]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *repository.UserService) GetUserByEmail(email string) (*repository.User, error) {
	login, exists := s.Emails[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	user, exists := s.Users[login]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *repository.UserService) GetUserByPhoneNumber(phone_number string) (*repository.User, error) {
	login, exists := s.PhoneNumbers[phone_number]
	if !exists {
		return nil, errors.New("user not found")
	}
	user, exists := s.Users[login]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
