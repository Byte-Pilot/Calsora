package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler interface {
	GetById(c *gin.Context)
	DeleteId(c *gin.Context)
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) GetById(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := h.service.GetById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"bday":  user.Bday.Format("2006-01-02"),
	})
}

func (h *userHandler) DeleteId(c *gin.Context) {
	userID := c.GetInt("user_id")
	err := h.service.DeleteId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
