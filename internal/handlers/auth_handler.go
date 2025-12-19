package handlers

import (
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type AuthHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type AuthHandler struct {
	service services.AuthServiceInterface
}

func NewAuthHandler(service services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (a *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Bday     string `json:"bday"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bday, err := time.Parse("2006-01-02", req.Bday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	access, refresh, err := a.service.Register(req.Email, req.Password, bday)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (a *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := a.service.Login(req.Email, req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (a *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := a.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
