package handlers

import (
	"Calsora/internal/httphelpers"
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	Logout(c *gin.Context)
	ChangePass(c *gin.Context)
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *authHandler {
	return &authHandler{
		service: service,
	}
}

func (a *authHandler) Register(c *gin.Context) {
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
	if httphelpers.IsWebClient(c) {
		httphelpers.SetAuthCookies(c, access, refresh)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (a *authHandler) Login(c *gin.Context) {
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
	if httphelpers.IsWebClient(c) {
		httphelpers.SetAuthCookies(c, access, refresh)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (a *authHandler) Refresh(c *gin.Context) {
	var refresh string
	var err error
	if httphelpers.IsWebClient(c) {
		refresh, err = c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
	} else {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		refresh = req.RefreshToken
	}

	access, refresh, err := a.service.Refresh(refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	if httphelpers.IsWebClient(c) {
		httphelpers.SetAuthCookies(c, access, refresh)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (a *authHandler) Logout(c *gin.Context) {
	var refresh string
	var err error

	if httphelpers.IsWebClient(c) {
		refresh, err = c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		httphelpers.ClearAuthCookies(c)
	} else {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		refresh = req.RefreshToken
	}
	err = a.service.Logout(refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (a *authHandler) ChangePass(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	access, refresh, err := a.service.ChangePass(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid input"})
		return
	}

	if httphelpers.IsWebClient(c) {
		httphelpers.SetAuthCookies(c, access, refresh)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
