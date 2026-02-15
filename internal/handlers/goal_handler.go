package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/services"
	"github.com/labstack/echo/v4"
)

type GoalHandler struct {
	GoalService *services.GoalService
	UserService *services.UserService
}

func NewGoalHandler(goalService *services.GoalService, userService *services.UserService) *GoalHandler {
	return &GoalHandler{
		GoalService: goalService,
		UserService: userService,
	}
}

func (h *GoalHandler) GetCurrentGoals(c echo.Context) error {
	// For simplicity, we'll assume the username is passed as a query param or auth token
	// In a real app, this would come from JWT middleware
	username := c.QueryParam("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username query parameter is required"})
	}

	ctx := c.Request().Context()
	user, err := h.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	goals, err := h.GoalService.GetUserGoals(ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) GenerateGoals(c echo.Context) error {
	// Manually trigger goal generation (for testing/admin)
	username := c.Param("username")
	ctx := c.Request().Context()

	user, err := h.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	if err := h.GoalService.GenerateWeeklyGoals(ctx, user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Goals generated successfully"})
}
