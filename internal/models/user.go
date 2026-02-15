package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User represents the comprehensive user data model
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// ==========================================
	// Authentication & Identity
	// ==========================================
	Username     string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"` // Never return password hash in JSON

	// ==========================================
	// Basic Profile Info (from getUserProfile)
	// ==========================================
	Name     string     `json:"name"`
	Avatar   string     `json:"avatar"`
	About    string     `gorm:"type:text" json:"about"`
	Birthday *time.Time `json:"birthday"`
	Country  string     `json:"country"`
	School   string     `json:"school"`
	Company  string     `json:"company"`

	// ==========================================
	// Social Links (from getUserProfile)
	// ==========================================
	GitHub   string         `json:"gitHub"`
	Twitter  string         `json:"twitter"`
	LinkedIn string         `json:"linkedIn"`
	Websites datatypes.JSON `gorm:"type:jsonb" json:"websites"` // Stores []string

	// ==========================================
	// LeetCode Stats (from getUserStats)
	// ==========================================
	Ranking           int `json:"ranking"`
	Reputation        int `json:"reputation"`
	ContributionPoint int `json:"contributionPoint"`

	// Problem Solving Counts
	TotalSolved  int `json:"totalSolved"`
	EasySolved   int `json:"easySolved"`
	MediumSolved int `json:"mediumSolved"`
	HardSolved   int `json:"hardSolved"`

	// ==========================================
	// Skills & Topics (from getUserSkills)
	// ==========================================
	// Stores the detailed list of topics and problem counts
	// Structure: { "fundamental": [...], "intermediate": [...], "advanced": [...] }
	SkillTags datatypes.JSON `gorm:"type:jsonb" json:"skillTags"`

	// ==========================================
	// Contest Performance (from getUserContest)
	// ==========================================
	ContestRating        float64 `json:"contestRating"`
	ContestGlobalRanking int     `json:"contestGlobalRanking"`
	ContestTopPercentage float64 `json:"contestTopPercentage"`
	TotalParticipants    int     `json:"totalParticipants"`
	ContestAttended      int     `json:"contestAttended"` // Count of contests attended

	// Badges (from getUserContest or getUserBadges)
	Badges datatypes.JSON `gorm:"type:jsonb" json:"badges"`

	// ==========================================
	// Activity & Consistency (from getUserCalendar)
	// ==========================================
	Streak          int            `json:"streak"`
	TotalActiveDays int            `json:"totalActiveDays"`
	ActiveYears     datatypes.JSON `gorm:"type:jsonb" json:"activeYears"` // []int

	// Submission Calendar (Heatmap)
	SubmissionCalendar datatypes.JSON `gorm:"type:jsonb" json:"submissionCalendar"`

	// ==========================================
	// App-Specific Calculated Score
	// ==========================================
	TotalScore     int            `json:"totalScore"`
	ScoreRank      string         `json:"scoreRank"`                        // Beginner, Intermediate, etc.
	ScoreBreakdown datatypes.JSON `gorm:"type:jsonb" json:"scoreBreakdown"` // Breakdown details
}

// TableName overrides the default table name if needed
func (User) TableName() string {
	return "users"
}
