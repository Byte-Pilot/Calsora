package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
)

type UserService interface {
	Create(user *models.User) error
	GetById(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	ChangePass(id int, newPass string) error
	DeleteId(id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) Create(user *models.User) error {
	return s.repo.Create(user)
}
func (s *userService) GetById(id int) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) ChangePass(id int, newPass string) error {
	return s.repo.ChangePass(id, newPass)
}
func (s *userService) DeleteId(id int) error {
	return s.repo.DeleteId(id)
}
