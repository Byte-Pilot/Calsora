package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MealHandler interface {
	AddMeal(c *gin.Context)
	DeleteMeal(c *gin.Context)
	GetDailyNutritionStats(c *gin.Context)
}

type mealHandler struct {
	service services.MealService
}

func NewMealHandler(svc services.MealService) *mealHandler {
	return &mealHandler{service: svc}
}

type addMealReq struct {
	Description string `json:"description"`
	PhotoURL    string `json:"photo_url"`
	UserID      int    `json:"user_id"`
}

func (mh *mealHandler) AddMeal(c *gin.Context) {
	var req addMealReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meal, err := mh.service.AddMeal(req.UserID, req.Description, req.PhotoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, meal)
}

func (mh *mealHandler) DeleteMeal(c *gin.Context) {
	userID := c.GetInt("user_id")
	err := mh.service.DeleteMeal(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meal not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (mh *mealHandler) GetDailyNutritionStats(c *gin.Context) {
	daysStr := c.Param("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid days"})
		return
	}
	stats, err := mh.service.GetDailyNutritionStats(days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, stats)
}
