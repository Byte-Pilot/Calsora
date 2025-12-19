package server

import (
	"Calsora/internal/handlers"
	"Calsora/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RunServer(userHandler handlers.UserHandlerInterface, authHandler handlers.AuthHandlerInterface, mealHandler handlers.MealHandlerInterface) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)

		protected := api.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/users/", userHandler.GetById)
			protected.DELETE("/users/delete/", userHandler.DeleteId)

			protected.POST("meals/add", mealHandler.AddMeal)
			protected.DELETE("meals/delete/", mealHandler.DeleteMeal)
		}
	}

	return r
}
