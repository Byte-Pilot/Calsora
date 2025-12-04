package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
)

type MealServiceInterface interface {
	AddMeal(userID uint, description string, photoURL string) (*models.Meals, error)
}

type MealService struct {
	repo repository.MealRepositoryInterface
}

func NewMealService(repo repository.MealRepositoryInterface) *MealService {
	return &MealService{repo: repo}
}

func (s *MealService) AddMeal(userID uint, description string, photoURL string) (*models.Meals, error) {
	// gpt
	var protein, carbs, fats float32 = 99.0, 99.0, 99.0
	var cal uint = 99
	name := "TestMeal"

	meal := &models.Meals{
		UserID:  userID,
		Name:    name,
		Cal:     cal,
		Protein: protein,
		Carbs:   carbs,
		Fats:    fats,
	}

	if err := s.repo.CreateMeal(meal); err != nil {
		return nil, err
	}
	return meal, nil
}
