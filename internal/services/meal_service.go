package services

import (
	"Calsora/internal/apperrors"
	"Calsora/internal/image"
	"Calsora/internal/intelligence/inference"
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

type MealService interface {
	AddMeal(userID int, description string, photoBytes []byte) (*models.Meals, []*models.MealItems, bool, error)
	EditMeal(userID, mealID int, description string, photoBytes []byte) (*models.Meals, []*models.MealItems, bool, error)
	UpdateMeal(meal *models.Meals, items []*models.MealItems) error
	GetDailyNutritionStats(userID, days int) ([]*models.Stats, error)
	DeleteMeal(userID, mealID int) error
}

type mealService struct {
	repo repository.MealRepository
	gpt  *inference.GPTClient
}

func NewMealService(repo repository.MealRepository, gpt *inference.GPTClient) *mealService {
	return &mealService{repo: repo, gpt: gpt}
}

func (s *mealService) AddMeal(userID int, description string, photoBytes []byte) (*models.Meals, []*models.MealItems, bool, error) {
	descRunes := []rune(description)
	if len(descRunes) > 1000 {
		description = string(descRunes[:1000])
	}

	var encodedPhoto string
	if photoBytes != nil {
		validPhoto, err := image.ImageProccesor(photoBytes, 1024)
		if err != nil {
			return nil, nil, false, err
		}
		encodedPhoto = base64.StdEncoding.EncodeToString(validPhoto)
	}

	result, err := s.gpt.AnalyzeMeal(description, encodedPhoto)
	log.Printf("AI response: %+v", result)
	if err != nil {
		return nil, nil, false, apperrors.NewCustomError("AI request failed", "К сожалению, обработать запрос не удалось")
	}

	meal := &models.Meals{
		UserID: userID,
		Name:   result.Name,
	}

	var items []*models.MealItems
	var avgConfidence float64
	for _, item := range result.Items {
		mealItem := &models.MealItems{
			Name:       item.Name,
			Grams:      item.Grams,
			Cal:        item.Calories,
			Protein:    item.Protein,
			Carbs:      item.Carbs,
			Fats:       item.Fats,
			Confidence: item.Confidence,
		}
		avgConfidence += item.Confidence
		items = append(items, mealItem)
	}
	if len(items) > 0 {
		avgConfidence /= float64(len(items))
	}
	lowConfidence := avgConfidence < 0.6

	if meal.ID, meal.CreatedAt, items, err = s.repo.CreateMeal(meal, items); err != nil {
		return nil, nil, false, err
	}

	return meal, items, lowConfidence, err
}

func (s *mealService) EditMeal(userID, mealID int, description string, photoBytes []byte) (*models.Meals, []*models.MealItems, bool, error) {
	descRunes := []rune(description)
	if len(descRunes) > 1000 {
		description = string(descRunes[:1000])
	}

	var encodedPhoto string
	if photoBytes != nil {
		validPhoto, err := image.ImageProccesor(photoBytes, 1024)
		if err != nil {
			return nil, nil, false, err
		}
		encodedPhoto = base64.StdEncoding.EncodeToString(validPhoto)
	}

	meal, items, err := s.repo.GetMealData(userID, mealID)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed DB request: %w", err)
	}
	meal.ID, meal.UserID = mealID, userID

	var promptItems []inference.MealAnalysisItems
	for _, item := range items {
		promptItems = append(promptItems, inference.MealAnalysisItems{
			Name:     item.Name,
			Grams:    item.Grams,
			Calories: item.Cal,
			Protein:  item.Protein,
			Carbs:    item.Carbs,
			Fats:     item.Fats,
		})
	}

	mealJSON, err := json.Marshal(meal.Name)
	itemsJSON, err := json.Marshal(promptItems)
	if err != nil {
		return nil, nil, false, apperrors.NewCustomError("failed marshal json", "К сожалению, обработать запрос не удалось")
	}

	respText := fmt.Sprintf("User context: %s\nCurrent meal name: %s\nCurrent meal items %s\n"+
		"Task: Update the meal data according to user instructions and return the final corrected values.",
		description, mealJSON, itemsJSON)

	result, err := s.gpt.AnalyzeMeal(respText, encodedPhoto)
	log.Printf("AI response: %+v", result)
	if err != nil {
		return nil, nil, false, apperrors.NewCustomError("AI request failed", "К сожалению, обработать запрос не удалось")
	}

	meal.Name = result.Name

	var itemsNew []*models.MealItems
	var avgConfidence float64
	for _, item := range result.Items {
		mealItem := &models.MealItems{
			Name:       item.Name,
			Grams:      item.Grams,
			Cal:        item.Calories,
			Protein:    item.Protein,
			Carbs:      item.Carbs,
			Fats:       item.Fats,
			Confidence: item.Confidence,
		}
		avgConfidence += item.Confidence
		itemsNew = append(itemsNew, mealItem)
	}
	if len(itemsNew) > 0 {
		avgConfidence /= float64(len(itemsNew))
	}
	lowConfidence := avgConfidence < 0.6

	if itemsNew, err = s.repo.EditMeal(meal, itemsNew); err != nil {
		return nil, nil, false, err
	}

	return meal, itemsNew, lowConfidence, err
}

func (s *mealService) UpdateMeal(meal *models.Meals, items []*models.MealItems) error {
	return s.repo.UpdateMeal(meal, items)
}

func (s *mealService) GetDailyNutritionStats(userID, days int) ([]*models.Stats, error) {
	stats, err := s.repo.GetDailyNutritionStats(userID, days)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *mealService) DeleteMeal(userID, mealID int) error {
	return s.repo.DeleteMeal(userID, mealID)
}
