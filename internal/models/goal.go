package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type GoalDefinition struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	Type                string `gorm:"not null" json:"type"` // e.g., 'SOLVE_PROBLEMS'
	DescriptionTemplate string `json:"description_template"`
	DifficultyLevel     string `json:"difficulty_level"` // 'BEGINNER', 'INTERMEDIATE', 'ADVANCED'
}

type WeeklyGoal struct {
	ID                  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID              uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User                User           `gorm:"foreignKey:UserID" json:"-"`
	WeekStartDate       time.Time      `gorm:"not null" json:"week_start_date"`
	GoalType            string         `gorm:"not null;default:'GENERATED'" json:"goal_type"`
	DifficultyBreakdown datatypes.JSON `json:"difficulty_breakdown"` // JSON: {easy: 3, medium: 4, hard: 1}
	SelectedProblems    datatypes.JSON `json:"selected_problems"`    // JSON: List of problem IDs/details
	FocusTopics         datatypes.JSON `json:"focus_topics"`         // JSON: ["DP", "Graph"]
	CompletionPercent   float64        `gorm:"default:0" json:"completion_percent"`
	Status              string         `gorm:"default:'PENDING'" json:"status"` // 'PENDING', 'COMPLETED', 'FAILED'
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}
