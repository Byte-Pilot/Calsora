package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MealHandlerInterface interface {
	AddMeal(c *gin.Context)
}

type mealHandler struct {
	service services.MealServiceInterface
}

func NewMealHandler(svc services.MealServiceInterface) *mealHandler {
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
