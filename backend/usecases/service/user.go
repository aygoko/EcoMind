package service

import (
	repository "github.com/aygoko/EcoMind/backend/domain"
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
