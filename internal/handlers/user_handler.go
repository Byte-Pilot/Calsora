package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserHandlerInterface interface {
	Register(c *gin.Context)
	GetById(c *gin.Context)
	DeleteId(c *gin.Context)
}

type userHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(service services.UserServiceInterface) *userHandler {
	return &userHandler{service: service}
}

type registerRequest struct {
	Email    string
	Password string
	Bday     string
}

func (h *userHandler) Register(c *gin.Context) {
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
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email})
}

func (h *userHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	user, err := h.service.GetById(uint(id))
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.service.DeleteId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
