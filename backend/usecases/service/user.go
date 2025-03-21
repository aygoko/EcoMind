package service

import (
	repository "github.com/aygoko/EcoMind/backend/domain"
	"github.com/aygoko/EcoMind/backend/repsitory/errors"
	"github.com/google/uuid"
)

type UserService struct {
	Repo repository.UserService
}

func NewUserService(repo repository.UserService) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) FindOrCreateUserByProvider(provider string, providerUserID string, email string) (*repository.User, error) {
	user, err := s.Repo.GetByProvider(provider, providerUserID)
	if err == nil {
		return user, nil
	}

	if email != "" {
		_, err := s.Repo.GetByEmail(email)
		if err == nil {
			return nil, errors.New("email already exists but not linked to this provider")
		}
	}

	newUser := &repository.User{
		ID:             GenerateUserID(),
		Provider:       provider,
		ProviderUserID: providerUserID,
		Email:          email,
	}

	err = s.Repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetByProvider(provider string, providerUserID string) (*repository.User, error) {
	return s.Repo.GetByProvider(provider, providerUserID)
}

func (s *UserService) Get(login string) (*repository.User, error) {
	return s.Repo.Get(login)
}

func (s *UserService) GetByEmail(email string) (*repository.User, error) {
	return s.Repo.GetByEmail(email)
}

func (s *UserService) GetByPhoneNumber(phone_number string) (*repository.User, error) {
	return s.Repo.GetByPhoneNumber(phone_number)
}

func GenerateUserID() string {
	return uuid.New().String()
}
