package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
)

type UserServiceInterface interface {
	GetById(id uint) (*models.User, error)
	DeleteId(id uint) error
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetById(id uint) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *UserService) DeleteId(id uint) error {
	return s.repo.DeleteId(id)
}
