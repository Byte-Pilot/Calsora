package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserHandlerInterface interface {
	Register(c *gin.Context)
	GetById(c *gin.Context)
}

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(service services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

type registerRequest struct {
	Email    string
	Password string
	Bday     string
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bday, err := time.Parse("2006-01-02", req.Bday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	user, err := h.service.Register(req.Email, req.Password, bday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email})
}



}
