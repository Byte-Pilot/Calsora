package server

import (
	"Calsora/internal/auth/middleware"
	"Calsora/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RunServer(userHandler handlers.UserHandler, authHandler handlers.AuthHandler, mealHandler handlers.MealHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)

		protected := api.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.POST("/auth/change-password", authHandler.ChangePass)

			protected.POST("/users/", userHandler.GetById)
			protected.DELETE("/users/delete/", userHandler.DeleteId)

			protected.POST("meals/add", mealHandler.AddMeal)
			protected.DELETE("meals/delete/", mealHandler.DeleteMeal)
		}
	}

	return r
}
