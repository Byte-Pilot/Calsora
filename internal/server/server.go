package server

import (
	"Calsora/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RunServer(userHandler handlers.UserHandlerInterface, mealHandler handlers.MealHandlerInterface) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/users/register", userHandler.Register)
		api.GET("/users/:id", userHandler.GetById)
		api.DELETE("/users/delete/:id", userHandler.DeleteId)

		api.POST("meals/add", mealHandler.AddMeal)
		api.DELETE("meals/delete/:id", mealHandler.DeleteMeal)
	}

	return r
}
