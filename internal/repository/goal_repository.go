package repository

import (
	"context"
	"time"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalRepository interface {
	CreateBatch(ctx context.Context, goals []models.WeeklyGoal) error
	GetWeeklyGoals(ctx context.Context, userID uuid.UUID, weekStart time.Time) ([]models.WeeklyGoal, error)
	CreateGoalDefinition(ctx context.Context, def *models.GoalDefinition) error
	GetGoalDefinitions(ctx context.Context) ([]models.GoalDefinition, error)
}

type goalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) CreateBatch(ctx context.Context, goals []models.WeeklyGoal) error {
	return r.db.WithContext(ctx).Create(&goals).Error
}

func (r *goalRepository) GetWeeklyGoals(ctx context.Context, userID uuid.UUID, weekStart time.Time) ([]models.WeeklyGoal, error) {
	var goals []models.WeeklyGoal
	// Assuming weekStart is the beginning of the week, we might want to filter by range or exact match
	// For simplicity, let's match exact date for now, or finding goals created within that week
	err := r.db.WithContext(ctx).Where("user_id = ? AND week_start_date = ?", userID, weekStart).Find(&goals).Error
	return goals, err
}

func (r *goalRepository) CreateGoalDefinition(ctx context.Context, def *models.GoalDefinition) error {
	return r.db.WithContext(ctx).Create(def).Error
}

func (r *goalRepository) GetGoalDefinitions(ctx context.Context) ([]models.GoalDefinition, error) {
	var defs []models.GoalDefinition
	err := r.db.WithContext(ctx).Find(&defs).Error
	return defs, err
}
