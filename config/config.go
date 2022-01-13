package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type DbConfig struct {
	Name     string
	User     string
	Password string
	Type     string
	Host     string
	Port     int
}

type Config struct {
	Db DbConfig
}

func New() *Config {
	load()

	return &Config{
		Db: DbConfig{
			Name:     getEnv("DB_NAME", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASS", ""),
			Type:     getEnv("DB_TYPE", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
		},
	}
}

// Возвращает стоковую переменную окружения
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Возвращает целочисленную переменную окружения
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Возвращает булиновую переменную окружения
/*func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}*/

// Возвращает переменную окружения в виде слайса строк
/*func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}*/

func load() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
