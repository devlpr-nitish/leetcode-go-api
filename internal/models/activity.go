package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ActivityType string    `json:"activity_type"` // 'PROBLEM_SOLVED', 'CONTEST_RANK'
	ReferenceID  string    `json:"reference_id"`
	Timestamp    time.Time `json:"timestamp"`
}
