package config

import (
	"os"

	"github.com/joho/godotenv"
)

type RapidAPIConfig struct {
	APIKey string
	Host   string
}

func GetRapidAPIConfig() *RapidAPIConfig {
	// Try to load .env from the config directory
	err := godotenv.Load(".env.example")
	if err != nil {
		panic("💥 error loading .env file")
	}

	return &RapidAPIConfig{
		APIKey: GetEnv("RAPIDAPI_KEY", "super_secret_key"),
		Host:   GetEnv("RAPIDAPI_HOST", "fakehostname.com"),
	}
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
