package services

import (
	"Calsora/internal/models"
	"Calsora/internal/repository"
)

type MealServiceInterface interface {
	AddMeal(userID int, description string, photoURL string) (*models.Meals, error)
	GetDailyNutritionStats(days int) (int, error)
	DeleteMeal(userID int) error
}

type MealService struct {
	repo repository.MealRepositoryInterface
}

func NewMealService(repo repository.MealRepositoryInterface) *MealService {
	return &MealService{repo: repo}
}

func (s *MealService) AddMeal(userID int, description string, photoURL string) (*models.Meals, error) {
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

func (s *MealService) GetDailyNutritionStats(days int) (int, error) {
	var test = 0
	return test, nil
}

func (s *MealService) DeleteMeal(mealID int) error {
	return s.repo.DeleteMeal(mealID)
}
