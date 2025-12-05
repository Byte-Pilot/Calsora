package services

import (
	"Calsora/internal/Error"
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"Calsora/internal/utils"
	"fmt"
	"time"
)

type UserServiceInterface interface {
	Register(email, password string, bday time.Time) (*models.User, error)
	GetById(id uint) (*models.User, error)
	DeleteId(id uint) error
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password string, bday time.Time) (*models.User, error) {

	if err := utils.ValidatePass(password); err != nil {
		return nil, fmt.Errorf("validatePass: ", err)
	}

	hashed, err := utils.HashPass(password)
	if err != nil {
		return nil, fmt.Errorf("hashPass: ", err)
	}

	var user = &models.User{
		Email:    email,
		Password: hashed,
		Bday:     bday,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, Error.NewCustomError("repo create: "+err.Error(), "Не удалось создать юзера")
	}
	return user, nil
}

func (s *UserService) GetById(id uint) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *UserService) DeleteId(id uint) error {
	return s.repo.DeleteId(id)
}
