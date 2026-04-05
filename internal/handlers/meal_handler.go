package handlers

import (
	"Calsora/internal/httphelpers"
	"Calsora/internal/models"
	"Calsora/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MealHandler interface {
	AddMeal(c *gin.Context)
	EditMeal(c *gin.Context)
	UpdateMeal(c *gin.Context)
	GetDailyNutritionStats(c *gin.Context)
	DeleteMeal(c *gin.Context)
}

type mealHandler struct {
	service services.MealService
}

func NewMealHandler(svc services.MealService) *mealHandler {
	return &mealHandler{service: svc}
}

type addMealReq struct {
	Description string `form:"description"`
	PhotoBytes  []byte
	UserID      int
}

type addMealResponse struct {
	Meal          *models.Meals       `json:"meal"`
	Items         []*models.MealItems `json:"items"`
	LowConfidence bool                `json:"low_confidence"`
}

type editMealReq struct {
	UserID      int    `form:"user_id"`
	MealID      int    `form:"meal_id"`
	Description string `form:"description"`
	PhotoBytes  []byte
}

func (mh *mealHandler) AddMeal(c *gin.Context) {
	var req addMealReq

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserID = c.GetInt("user_id")

	if _, err := c.FormFile("image"); err == nil {
		file, err := httphelpers.UploadImage(c, 10<<20)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.PhotoBytes = file
	}

	meal, items, lowConfidence, err := mh.service.AddMeal(req.UserID, req.Description, req.PhotoBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, addMealResponse{Meal: meal, Items: items, LowConfidence: lowConfidence})

	/* TEST
	c.Data(http.StatusOK, "image/jpeg", photo)
	*/
}

func (mh *mealHandler) EditMeal(c *gin.Context) {
	var editReq editMealReq

	if err := c.ShouldBind(&editReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	editReq.UserID = c.GetInt("user_id")

	if _, err := c.FormFile("image"); err == nil {
		file, err := httphelpers.UploadImage(c, 10<<20)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		editReq.PhotoBytes = file
	}

	meal, items, lowConfidence, err := mh.service.EditMeal(editReq.UserID, editReq.MealID, editReq.Description, editReq.PhotoBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, addMealResponse{Meal: meal, Items: items, LowConfidence: lowConfidence})

}

func (mh *mealHandler) UpdateMeal(c *gin.Context) {
	userID := c.GetInt("user_id")
	var updateMealReq struct {
		Meal  *models.Meals       `json:"meal"`
		Items []*models.MealItems `json:"items"`
	}
	if err := c.ShouldBindJSON(&updateMealReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid JSON request": err.Error()})
		return
	}
	if updateMealReq.Meal == nil || updateMealReq.Items == nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid JSON request": errors.New("meal and items must be provided")})
		return
	}

	updateMealReq.Meal.UserID = userID

	if err := mh.service.UpdateMeal(updateMealReq.Meal, updateMealReq.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (mh *mealHandler) GetDailyNutritionStats(c *gin.Context) {
	userID := c.GetInt("user_id")
	daysStr := c.DefaultQuery("days", "1")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid days", "days value": err})
		return
	}
	if days > 365 {
		days = 365
	}

	stats, err := mh.service.GetDailyNutritionStats(userID, days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (mh *mealHandler) DeleteMeal(c *gin.Context) {
	userID := c.GetInt("user_id")
	mealIDStr := c.Param("id")
	mealID, err := strconv.Atoi(mealIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meal ID"})
		return
	}

	err = mh.service.DeleteMeal(userID, mealID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
