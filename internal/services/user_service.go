package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
)

type UserService interface {
	GetById(id int) (*models.User, error)
	DeleteId(id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetById(id int) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *userService) DeleteId(id int) error {
	return s.repo.DeleteId(id)
}
