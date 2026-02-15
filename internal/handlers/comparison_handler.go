package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/services"
	"github.com/labstack/echo/v4"
)

type ComparisonHandler struct {
	service *services.ComparisonService
}

func NewComparisonHandler(service *services.ComparisonService) *ComparisonHandler {
	return &ComparisonHandler{service: service}
}

type CompareRequest struct {
	User1Name string `json:"user1_name"`
	User2Name string `json:"user2_name"`
	User1Data string `json:"user1_data"`
	User2Data string `json:"user2_data"`
}

func (h *ComparisonHandler) CompareUsers(c echo.Context) error {
	var req CompareRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.User1Data == "" || req.User2Data == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Both user1_data and user2_data are required"})
	}

	// Default names if not provided
	if req.User1Name == "" {
		req.User1Name = "User 1"
	}
	if req.User2Name == "" {
		req.User2Name = "User 2"
	}

	result, err := h.service.CompareUsers(c.Request().Context(), req.User1Name, req.User1Data, req.User2Name, req.User2Data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
