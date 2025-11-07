package server

import (
	"Calsora/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RunServer(userHandler handlers.UserHandlerInterface) error {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.GET("/users/:id", userHandler.GetById)
	}

	return r.Run(":8080")
}
