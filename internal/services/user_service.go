package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"time"
)

type UserServiceInterface interface {
	Register(email, password string, bday time.Time) (*models.User, error)
	GetById(id int) (*models.User, error)
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password string, bday time.Time) (*models.User, error) {
	var user = &models.User{
		Email:     email,
		Password:  password,
		Bday:      bday,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
