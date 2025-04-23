
package service

import (
    "github.com/aygoko/EcoMInd/backend/domain"
)

type UserService struct {
    Repo domain.UserService 
}

func NewUserService(repo domain.UserService) *UserService {
    return &UserService{
        Repo: repo,
    }
}

func (s *UserService) Get(login string) (*domain.User, error) {
    return s.Repo.Get(login)
}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
    return s.Repo.GetByEmail(email)
}

func (s *UserService) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
    return s.Repo.GetByPhoneNumber(phoneNumber)
}