package services

import (
	"Calsora/internal/apperrors"
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
	ChangePass(userID int, oldPassword, newPassword string) (string, string, error)
}

type authService struct {
	authRepo    repository.AuthRepository
	userService UserService
	subService  SubService
}

func NewAuthService(a repository.AuthRepository, u UserService, s SubService) *authService {
	return &authService{
		authRepo:    a,
		userService: u,
		subService:  s,
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
	if err := s.userService.Create(user); err != nil {
		return "", "", apperrors.NewCustomError("repo create: "+err.Error(), "Не удалось создать пользователя")
	}

	userSub, err := s.subService.Create(user.ID)
	if err != nil {
		s.userService.DeleteId(user.ID)
		return "", "", apperrors.NewCustomError("repo create: "+err.Error(), "Не удалось создать подписку")
	}

	access, err := jwt.GenerateAccessToken(user.ID, userSub.Plan, userSub.Status)
	if err != nil {
		s.userService.DeleteId(user.ID)
		return "", "", apperrors.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}
	refresh, exp, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		s.userService.DeleteId(user.ID)
		return "", "", apperrors.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
	}

	err = s.authRepo.SaveRefreshToken(user.ID, refresh, exp)
	if err != nil {
		s.userService.DeleteId(user.ID)
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	var user = &models.User{}
	user, err := s.userService.GetByEmail(email)
	if err != nil || !hash.CheckPass(password, user.Password) {
		return "", "", apperrors.NewCustomError("invalid login data: ", "Неверный логин или пароль")
	}

	userSub, err := s.subService.GetSubDataByID(user.ID)
	if err != nil {
		return "", "", apperrors.NewCustomError("GetSubDataByID failed: ", "Не удалось получить данные о подписке")
	}

	access, err := jwt.GenerateAccessToken(user.ID, userSub.Plan, userSub.Status)
	if err != nil {
		return "", "", apperrors.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}
	refresh, exp, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", apperrors.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
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
		return "", "", apperrors.NewCustomError("invalid token", "Сессия истекла, войдите снова")
	}

	userSub, err := s.subService.GetSubDataByID(userID)
	if err != nil {
		return "", "", apperrors.NewCustomError("GetSubDataByID failed: "+err.Error(), "Не удалось получить данные о подписке")
	}

	access, err := jwt.GenerateAccessToken(userID, userSub.Plan, userSub.Status)
	if err != nil {
		return "", "", apperrors.NewCustomError("generate access token: "+err.Error(), "Что-то пошло не так")
	}

	refreshNEW, exp, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", apperrors.NewCustomError("generate refresh token: "+err.Error(), "Что-то пошло не так")
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
		return apperrors.NewCustomError("delete refresh token: "+err.Error(), "Что-то пошло не так")
	}
	return nil
}

func (s *authService) ChangePass(userID int, oldPassword, newPassword string) (string, string, error) {

	user, err := s.userService.GetById(userID)
	if err != nil {
		return "", "", fmt.Errorf("GetByID: %w", err)
	}
	if !hash.CheckPass(oldPassword, user.Password) {
		return "", "", apperrors.NewCustomError("oldPassword invalid: "+err.Error(), "Актуальный пароль введен неверно")
	}
	if err := validator.ValidatePass(newPassword); err != nil {
		return "", "", apperrors.NewCustomError("newPassword invalid: "+err.Error(), "Новый пароль не соответствует требованиям")
	}

	newHashed, err := hash.HashPass(newPassword)
	if err != nil {
		return "", "", fmt.Errorf("hashPass: %w", err)
	}
	err = s.userService.ChangePass(userID, newHashed)
	if err != nil {
		return "", "", apperrors.NewCustomError("update password: "+err.Error(), "Неудалось обновить пароль")
	}

	err = s.authRepo.DeleteAllRefreshTokens(userID)
	if err != nil {
		return "", "", apperrors.NewCustomError("delete all refresh tokens: "+err.Error(), "Что-то пошло не так")
	}

	userSub, err := s.subService.GetSubDataByID(userID)
	if err != nil {
		return "", "", apperrors.NewCustomError("GetSubDataByID failed: "+err.Error(), "Войдите в учетную запись повторно")
	}

	newAccess, err := jwt.GenerateAccessToken(userID, userSub.Plan, userSub.Status)
	if err != nil {
		return "", "", apperrors.NewCustomError("generate access token: "+err.Error(), "Войдите в учетную запись повторно")
	}

	newRefresh, exp, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", apperrors.NewCustomError("generate refresh token: "+err.Error(), "Войдите в учетную запись повторно")
	}
	err = s.authRepo.SaveRefreshToken(userID, newRefresh, exp)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}
