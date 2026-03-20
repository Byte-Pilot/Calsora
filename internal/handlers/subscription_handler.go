package handlers

import (
	"Calsora/internal/httphelpers"
	"Calsora/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubHandler interface {
	CreateTrial(c *gin.Context)
	CreatePremium(c *gin.Context)
}

type subHandler struct {
	service services.SubService
}

func NewSubHandler(service services.SubService) *subHandler { return &subHandler{service: service} }

func (h *subHandler) CreateTrial(c *gin.Context) {

}

func (h *subHandler) CreatePremium(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Promo string `json:"promo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePremium(userID, req.Promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if httphelpers.IsWebClient(c) {
		httphelpers.ClearAuthCookiesAccess(c)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "access_token": nil})
}
