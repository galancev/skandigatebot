package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Host     string
	User     string
	Password string
}

func New() *Config {
	load()

	return &Config{
		Host:     getEnv("PACS_HOST", "localhost"),
		User:     getEnv("PACS_USER", ""),
		Password: getEnv("PACS_PASS", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func load() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
