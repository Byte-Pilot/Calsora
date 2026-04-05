package services

import (
	"Calsora/internal/intelligence/nutrition"
	"Calsora/internal/models"
	"Calsora/internal/repository"
	"time"
)

type UserProfileService interface {
	GetDailyIntake(profile *models.UserProfile, goal *models.UserGoal) (*models.NutritionTarget, error)
}

type userProfileService struct {
	repo        repository.UserProfileRepository
	userService UserService
}

func NewUserProfileService(repo repository.UserProfileRepository, userService UserService) *userProfileService {
	return &userProfileService{repo: repo, userService: userService}
}

func (s *userProfileService) GetDailyIntake(profile *models.UserProfile, goal *models.UserGoal) (*models.NutritionTarget, error) {
	userData, err := s.userService.GetById(profile.UserID)
	if err != nil {
		return nil, err
	}

	age := time.Now().Year() - userData.Bday.Year()
	if time.Now().YearDay() < userData.Bday.YearDay() {
		age--
	}

	dailyIntake := nutrition.CalculateNorm(profile, goal, age)
	dailyIntake.UserID = profile.UserID

	err = s.repo.GetDailyIntake(profile, goal, dailyIntake)
	if err != nil {
		return nil, err
	}

	return dailyIntake, nil
}
