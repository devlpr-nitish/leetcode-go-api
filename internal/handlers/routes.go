package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, userHandler *UserHandler, goalHandler *GoalHandler, authHandler *AuthHandler, comparisonHandler *ComparisonHandler) {
	api := e.Group("/api/v1")

	// Auth Routes
	api.POST("/auth/signup", authHandler.Signup)
	api.POST("/auth/login", authHandler.Login)

	// User Routes
	api.GET("/users/:username", userHandler.GetUser)
	api.POST("/users/:username/sync", userHandler.SyncUser)

	// Goal Routes
	api.GET("/goals/current", goalHandler.GetCurrentGoals)
	api.POST("/users/:username/goals/generate", goalHandler.GenerateGoals)

	// Comparison Routes
	api.POST("/compare", comparisonHandler.CompareUsers)

	// Health Check
	e.GET("/health", HealthCheck)
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
