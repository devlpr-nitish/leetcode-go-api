package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	LeetCodeAPI string
	JWTSecret   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DB_URL", "host=localhost user=postgres password=postgres dbname=leetcode_tracker port=5432 sslmode=disable"),
		LeetCodeAPI: getEnv("LEETCODE_API_URL", "https://leetcode.com/graphql"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-me"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
