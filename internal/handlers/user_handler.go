package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandlerInterface interface {
	GetById(c *gin.Context)
	DeleteId(c *gin.Context)
}

type userHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(service services.UserServiceInterface) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) GetById(c *gin.Context) {
	/*
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		user, err := h.service.GetById(uint(id))
	*/
	userID := c.GetInt("user_id")
	user, err := h.service.GetById(uint(userID))
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
	/*
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		err = h.service.DeleteId(uint(id))
	*/
	userID := c.GetInt("user_id")
	err := h.service.DeleteId(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
