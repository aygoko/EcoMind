package service

import (
    "github.com/aygoko/EcoMInd/backend/domain"
)

// UserService implements business logic for user operations.
type UserService struct {
    Repo domain.UserService
}

// NewUserService creates a new user service instance.
// Panics if the provided repository is nil.
func NewUserService(repo domain.UserService) *UserService {
    if repo == nil {
        panic("repository must not be nil")
    }
    return &UserService{
        Repo: repo,
    }
}

// Get retrieves a user by login.
func (s *UserService) Get(login string) (*domain.User, error) {
    return s.Repo.Get(login)
}

// GetByEmail retrieves a user by email.
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
    return s.Repo.GetByEmail(email)
}

// GetByPhoneNumber retrieves a user by phone number.
func (s *UserService) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
    return s.Repo.GetByPhoneNumber(phoneNumber)
}