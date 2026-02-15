package main

import (
	"log"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/config"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	repository.InitDB(cfg)

	// Drop legacy columns
	err := repository.DB.Exec(`
		ALTER TABLE weekly_goals 
		DROP COLUMN IF EXISTS target_value,
		DROP COLUMN IF EXISTS current_value,
		DROP COLUMN IF EXISTS details;
	`).Error

	if err != nil {
		log.Fatalf("Failed to drop columns: %v", err)
	}

	log.Println("Successfully dropped legacy columns from weekly_goals")
}
