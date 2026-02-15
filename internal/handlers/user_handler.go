package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/services"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	username := c.Param("username")
	ctx := c.Request().Context()

	user, err := h.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SyncUser(c echo.Context) error {
	username := c.Param("username")
	ctx := c.Request().Context()

	user, err := h.UserService.SyncUser(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
