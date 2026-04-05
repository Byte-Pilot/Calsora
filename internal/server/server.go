package server

import (
	"Calsora/internal/auth/middleware"
	"Calsora/internal/handlers"
	"Calsora/pkg/ratelimiter"
	"github.com/gin-gonic/gin"
)

func RunServer(userHandler handlers.UserHandler, uProfileHandler handlers.UserProfileHandler, subHandler handlers.SubHandler, authHandler handlers.AuthHandler, mealHandler handlers.MealHandler, rateLimiterService ratelimiter.Limiter) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.IPRateLimitMiddleware(rateLimiterService))
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth())
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.POST("/auth/logout-all", authHandler.LogoutAllSessions)
			protected.POST("/auth/change-password", authHandler.ChangePass)

			protected.POST("/subscription", subHandler.CreatePremium)

			protected.GET("/users", userHandler.GetById)
			protected.DELETE("/users/delete", userHandler.DeleteId)

			protected.POST("/profile", uProfileHandler.GetDailyIntake)

			premium := protected.Group("")
			premium.Use(middleware.RequireActiveSubscription())
			{
				premium.POST("/meals/add", mealHandler.AddMeal)
				premium.POST("/meals/edit", mealHandler.EditMeal)
			}

			protected.PATCH("/meals/update", mealHandler.UpdateMeal)
			protected.GET("/meals/stats", mealHandler.GetDailyNutritionStats)
			protected.DELETE("/meals/delete/:id", mealHandler.DeleteMeal)
		}
	}

	return r
}
