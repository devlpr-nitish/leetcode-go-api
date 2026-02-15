package repository

import (
	"log"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/config"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")

	// Auto Migrate
	err = DB.AutoMigrate(
		&models.User{},
		&models.GoalDefinition{},
		&models.WeeklyGoal{},
		&models.ActivityLog{},
		&models.UserComparison{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migration completed")
}
