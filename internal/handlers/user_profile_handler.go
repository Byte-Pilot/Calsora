package handlers

import (
	"Calsora/internal/models"
	"Calsora/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserProfileHandler interface {
	GetDailyIntake(c *gin.Context)
}

type userProfileHandler struct {
	service services.UserProfileService
}

func NewUserProfileHandler(service services.UserProfileService) *userProfileHandler {
	return &userProfileHandler{service: service}
}

type getDailyIntakeRequest struct {
	Profile *models.UserProfile `json:"profile"`
	Goal    *models.UserGoal    `json:"goal"`
}

func (h *userProfileHandler) GetDailyIntake(c *gin.Context) {
	userID := c.GetInt("user_id")
	var req getDailyIntakeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Profile == nil || req.Goal == nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid JSON request": errors.New("profile and goal must be provided")})
		return
	}
	req.Profile.UserID = userID
	req.Goal.UserID = userID

	dailyIntake, err := h.service.GetDailyIntake(req.Profile, req.Goal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось рассчитать норму", "Возникла ошибки": err})
		return
	}

	c.JSON(http.StatusOK, dailyIntake)
}
