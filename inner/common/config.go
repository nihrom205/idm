package common

import (
	"github.com/joho/godotenv"
	"os"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName string `validate:"required"`
	DSN          string `validate:"required"`
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) Config {
	_ = godotenv.Load(envFile)

	return Config{
		DbDriverName: os.Getenv("DB_DRIVER_NAME"),
		DSN:          os.Getenv("DB_DSN"),
	}
}
