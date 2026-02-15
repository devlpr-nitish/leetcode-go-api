package main

import (
	"context"
	"log"

	lcResult "github.com/devlpr-nitish/leetcode-tracker-backend/internal/clients/leetcode"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/config"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/handlers"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/repository"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/services"
	"github.com/devlpr-nitish/leetcode-tracker-backend/pkg/leetcode"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/genai"
)

func main() {

	cfg := config.LoadConfig()

	repository.InitDB(cfg)

	// Legacy Client for Problems
	leetcodeClient := leetcode.NewClient(cfg.LeetCodeAPI)

	// New Client for User Profile
	userLeetCodeClient := lcResult.NewClient()

	userRepo := repository.NewUserRepository(repository.DB)
	goalRepo := repository.NewGoalRepository(repository.DB)
	problemRepo := repository.NewProblemRepository(leetcodeClient)

	userService := services.NewUserService(userRepo, userLeetCodeClient)
	goalService := services.NewGoalService(userRepo, goalRepo, problemRepo)
	authService := services.NewAuthService(userRepo, cfg)

	userHandler := handlers.NewUserHandler(userService)
	goalHandler := handlers.NewGoalHandler(goalService, userService)
	authHandler := handlers.NewAuthHandler(authService)

	// GenAI Client
	genaiClient, err := genai.NewClient(context.Background(), nil)
	if err != nil {
		log.Printf("Warning: Failed to create GenAI client: %v", err)
	}

	comparisonService := services.NewComparisonService(genaiClient)
	comparisonHandler := handlers.NewComparisonHandler(comparisonService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	handlers.RegisterRoutes(e, userHandler, goalHandler, authHandler, comparisonHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		e.Logger.Fatal(err)
	}
}
