package repository

import (
	"encoding/json"
	"time"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/models"
	"gorm.io/datatypes"
)

type ComparisonRepository struct{}

func NewComparisonRepository() *ComparisonRepository {
	return &ComparisonRepository{}
}

// GetValidComparison searches for a recent comparison between two users (in either order)
func (r *ComparisonRepository) GetValidComparison(user1, user2 string, duration time.Duration) (*models.UserComparison, error) {
	var comparison models.UserComparison

	// Check for both (user1, user2) and (user2, user1)
	err := DB.Where(
		"((user1_name = ? AND user2_name = ?) OR (user1_name = ? AND user2_name = ?)) AND created_at > ?",
		user1, user2, user2, user1, time.Now().Add(-duration),
	).First(&comparison).Error

	if err != nil {
		return nil, err
	}
	return &comparison, nil
}

// SaveComparison saves a new comparison result to the database
func (r *ComparisonRepository) SaveComparison(user1, user2 string, result interface{}) error {
	// Convert result to JSON
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	comparison := models.UserComparison{
		User1Name: user1,
		User2Name: user2,
		Result:    datatypes.JSON(jsonBytes),
	}

	return DB.Create(&comparison).Error
}
