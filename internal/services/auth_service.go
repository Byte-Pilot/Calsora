package services

import (
	"Calsora/internal/Error"
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"Calsora/internal/utils"
	"fmt"
	"log"
	"time"
)

type AuthServiceInterface interface {
	Register(email, password string, bday time.Time) (string, string, error)
	Login(email, password string) (string, string, error)
	Refresh(refresh string) (string, string, error)
}

type AuthService struct {
	authRepo repository.AuthRepositoryInterface
	userRepo repository.UserRepositoryInterface
	subRepo  repository.SubscriptionsRepositoryInterface
}

func NewAuthService(a repository.AuthRepositoryInterface, u repository.UserRepositoryInterface, s repository.SubscriptionsRepositoryInterface) *AuthService {
	return &AuthService{
		authRepo: a,
		userRepo: u,
		subRepo:  s,
	}
}

func (s *AuthService) Register(email, password string, bday time.Time) (string, string, error) {
	if err := utils.ValidatePass(password); err != nil {
		return "", "", fmt.Errorf("validatePass: %w", err)
	}

	hashed, err := utils.HashPass(password)
	if err != nil {
		return "", "", fmt.Errorf("hashPass: %w", err)
	}

	var user = &models.User{
		Email:    email,
		Password: hashed,
		Bday:     bday,
	}
	if err := s.userRepo.Create(user); err != nil {
		return "", "", Error.NewCustomError("repo create: "+err.Error(), "Не удалось создать пользователя")
	}

	userSub, err := s.subRepo.Create(user.ID)
	if err != nil {
		return "", "", Error.NewCustomError("repo create: "+err.Error(), "Не удалось создать подписку")
	}

	access, err := utils.GenerateAccessToken(user.ID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}
	refresh, exp, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", Error.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
	}

	err = s.authRepo.SaveRefreshToken(user.ID, refresh, exp)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	var user = &models.User{}
	user, err := s.userRepo.GetByEmail(email)
	if err != nil || !utils.CheckPass(password, user.Password) {
		return "", "", Error.NewCustomError("invalid login data: "+err.Error(), "Неверный логин или пароль")
	}

	userSub, err := s.subRepo.GetSubDataByID(user.ID)
	if err != nil {
		return "", "", Error.NewCustomError("GetSubDataByID failed: "+err.Error(), "Не удалось получить данные о подписке")
	}

	access, err := utils.GenerateAccessToken(user.ID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}
	refresh, exp, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", Error.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
	}

	err = s.authRepo.SaveRefreshToken(user.ID, refresh, exp)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) Refresh(refresh string) (string, string, error) {
	userID, exp, err := s.authRepo.GetRefreshToken(refresh)
	if err != nil {
		return "", "", err
	}
	if exp.Before(time.Now()) {
		s.authRepo.DeleteRefreshToken(refresh)
		return "", "", Error.NewCustomError("invalid token", "Сессия истекла, войдите снова")
	}

	userSub, err := s.subRepo.GetSubDataByID(userID)
	if err != nil {
		return "", "", Error.NewCustomError("GetSubDataByID failed: "+err.Error(), "Не удалось получить данные о подписке")
	}
	access, err := utils.GenerateAccessToken(userID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}

	refreshNEW, exp, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", Error.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
	}

	err = s.authRepo.SaveRefreshToken(userID, refreshNEW, exp)
	if err != nil {
		return "", "", err
	}

	err = s.authRepo.DeleteRefreshToken(refresh)
	if err != nil {
		log.Println(err)
	}

	return access, refreshNEW, nil
}
