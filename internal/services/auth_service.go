package services

import (
	"Calsora/internal/Error"
	"Calsora/internal/auth/hash"
	"Calsora/internal/auth/jwt"
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"Calsora/internal/validator"
	"fmt"
	"log"
	"time"
)

type AuthService interface {
	Register(email, password string, bday time.Time) (string, string, error)
	Login(email, password string) (string, string, error)
	Refresh(refresh string) (string, string, error)
	Logout(refresh string) error
	ChangePass(useriID int, oldPassword, newPassword string) (string, string, error)
}

type authService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	subRepo  repository.SubscriptionsRepository
}

func NewAuthService(a repository.AuthRepository, u repository.UserRepository, s repository.SubscriptionsRepository) *authService {
	return &authService{
		authRepo: a,
		userRepo: u,
		subRepo:  s,
	}
}

func (s *authService) Register(email, password string, bday time.Time) (string, string, error) {
	if err := validator.ValidatePass(password); err != nil {
		return "", "", fmt.Errorf("validatePass: %w", err)
	}

	hashed, err := hash.HashPass(password)
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
		s.userRepo.DeleteId(user.ID)
		return "", "", Error.NewCustomError("repo create: "+err.Error(), "Не удалось создать подписку")
	}

	access, err := jwt.GenerateAccessToken(user.ID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		s.userRepo.DeleteId(user.ID)
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}
	refresh, exp, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		s.userRepo.DeleteId(user.ID)
		return "", "", Error.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
	}

	err = s.authRepo.SaveRefreshToken(user.ID, refresh, exp)
	if err != nil {
		s.userRepo.DeleteId(user.ID)
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	var user = &models.User{}
	user, err := s.userRepo.GetByEmail(email)
	if err != nil || !hash.CheckPass(password, user.Password) {
		return "", "", Error.NewCustomError("invalid login data: "+err.Error(), "Неверный логин или пароль")
	}

	userSub, err := s.subRepo.GetSubDataByID(user.ID)
	if err != nil {
		return "", "", Error.NewCustomError("GetSubDataByID failed: "+err.Error(), "Не удалось получить данные о подписке")
	}

	access, err := jwt.GenerateAccessToken(user.ID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}
	refresh, exp, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", Error.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
	}

	err = s.authRepo.SaveRefreshToken(user.ID, refresh, exp)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) Refresh(refresh string) (string, string, error) {
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
	access, err := jwt.GenerateAccessToken(userID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}

	refreshNEW, exp, err := jwt.GenerateRefreshToken(userID)
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
		return "", "", err
	}

	return access, refreshNEW, nil
}

func (s *authService) Logout(refresh string) error {
	err := s.authRepo.DeleteRefreshToken(refresh)
	if err != nil {
		return Error.NewCustomError("delete refresh token: "+err.Error(), "Что-то пошло не так")
	}
	return nil
}

func (s *authService) ChangePass(userID int, oldPassword, newPassword string) (string, string, error) {

	user, err := s.userRepo.GetById(userID)
	if err != nil {
		return "", "", fmt.Errorf("GetByID: %w", err)
	}
	if !hash.CheckPass(oldPassword, user.Password) {
		return "", "", Error.NewCustomError("oldPassword invalid: "+err.Error(), "Актуальный пароль введен неверно")
	}
	if err := validator.ValidatePass(newPassword); err != nil {
		return "", "", Error.NewCustomError("newPassword invalid: "+err.Error(), "Новый пароль не соответствует требованиям")
	}

	newHashed, err := hash.HashPass(newPassword)
	if err != nil {
		return "", "", fmt.Errorf("hashPass: %w", err)
	}
	err = s.userRepo.ChangePass(userID, newHashed)
	if err != nil {
		return "", "", Error.NewCustomError("update password: "+err.Error(), "Неудалось обновить пароль")
	}

	err = s.authRepo.DeleteAllRefreshTokens(userID)
	if err != nil {
		return "", "", Error.NewCustomError("delete all refresh tokens: "+err.Error(), "Что-то пошло не так")
	}
	/*
		токены должны быть удалены сразу после смены пароля
	*/

	userSub, err := s.subRepo.GetSubDataByID(userID)
	if err != nil {
		return "", "", Error.NewCustomError("GetSubDataByID failed: "+err.Error(), "Войдите в учетную запись повторно")
	}
	newAccess, err := jwt.GenerateAccessToken(userID, userSub.Plan, userSub.ExpiresAt)
	if err != nil {
		return "", "", Error.NewCustomError("generate access token: "+err.Error(), "Войдите в учетную запись повторно")
	}

	newRefresh, exp, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", Error.NewCustomError("generate refresh token: "+err.Error(), "Войдите в учетную запись повторно")
	}
	err = s.authRepo.SaveRefreshToken(userID, newRefresh, exp)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}
