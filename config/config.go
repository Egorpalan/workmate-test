package config

import (
	"fmt"
	"os"

	"github.com/Egorpalan/workmate-test/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Info("Warning: .env file not found, using environment variables")
	}

	dbConfig := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "tasks_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	serverConfig := ServerConfig{
		Port: getEnv("SERVER_PORT", "8080"),
	}

	return &Config{
		DB:     dbConfig,
		Server: serverConfig,
	}, nil
}

// getEnv получает значение из переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN возвращает строку подключения к базе данных
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
