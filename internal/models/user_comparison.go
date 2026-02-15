package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// UserComparison stores the result of a comparison between two users
type UserComparison struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User1Name string         `gorm:"index;not null" json:"user1_name"`
	User2Name string         `gorm:"index;not null" json:"user2_name"`
	Result    datatypes.JSON `gorm:"type:jsonb;not null" json:"result"` // Stores the JSON result from AI
}

// TableName overrides the default table name
func (UserComparison) TableName() string {
	return "user_comparisons"
}
