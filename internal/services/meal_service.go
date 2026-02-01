package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
)

type MealService interface {
	AddMeal(userID int, description string, photoURL string) (*models.Meals, error)
	GetDailyNutritionStats(days int) (int, error)
	DeleteMeal(userID int) error
}

type mealService struct {
	repo repository.MealRepository
}

func NewMealService(repo repository.MealRepository) *mealService {
	return &mealService{repo: repo}
}

func (s *mealService) AddMeal(userID int, description string, photoURL string) (*models.Meals, error) {
	// gpt
	description = ""
	photoURL = ""
	var protein, carbs, fats float32 = 99.0, 99.0, 99.0
	var cal int = 99
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

func (s *mealService) GetDailyNutritionStats(days int) (int, error) {
	var test = 0
	return test, nil
}

func (s *mealService) DeleteMeal(mealID int) error {
	return s.repo.DeleteMeal(mealID)
}
