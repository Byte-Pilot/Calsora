package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type SubService interface {
	Create(userID int) (*models.Subscriptions, error)
	CreateTrial(c *gin.Context) error
	CreatePremium(userID int, promo string) error
	GetSubDataByID(userID int) (*models.Subscriptions, error)
}

type subService struct {
	repo repository.SubscriptionsRepository
}

func NewSubService(repo repository.SubscriptionsRepository) *subService {
	return &subService{repo: repo}
}

func (s *subService) Create(userID int) (*models.Subscriptions, error) {
	return s.repo.Create(userID)
}

func (s *subService) CreateTrial(c *gin.Context) error {
	return nil
}

func (s *subService) CreatePremium(userID int, promo string) error {
	plan := "premium"
	exp := time.Now().Add(24 * 30 * time.Hour)
	promoKey := os.Getenv("PROMO")

	if promo == "" {
		return fmt.Errorf("invalid promo")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(promoKey), []byte(promo)); err == nil {
		if err := s.repo.CreatePremium(userID, plan, exp); err != nil {
			return errors.New("invalid promo")
		}
	} else {
		return fmt.Errorf("invalid promo %w", err)
	}

	return nil
}

func (s *subService) GetSubDataByID(userID int) (*models.Subscriptions, error) {
	userSub, err := s.repo.GetSubDataByID(userID)
	if err != nil {
		return nil, err
	}
	if userSub.Plan != "free" && userSub.Status == "active" && userSub.ExpiresAt.Before(time.Now()) {
		if err := s.repo.MarkExpired(userSub.ID); err != nil {
			log.Println("failed to mark expired", err)
		}
		userSub.Status = "expired"
	}

	return userSub, nil
}
