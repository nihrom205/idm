package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName string `validate:"required"`
	DSN          string `validate:"required"`
	AppName      string `validate:"required"`
	AppVersion   string `validate:"required"`
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) Config {
	// если нет файла, то залогируем это и попробуем получить конфиг из переменных окружения
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	cfg := Config{
		DbDriverName: os.Getenv("DB_DRIVER_NAME"),
		DSN:          os.Getenv("DB_DSN"),
		AppName:      os.Getenv("APP_NAME"),
		AppVersion:   os.Getenv("APP_VERSION"),
	}

	err = validator.New().Struct(&cfg)
	if err != nil {
		var validatorErr validator.ValidationErrors
		if errors.As(err, &validatorErr) {
			// если конфиг не прошел валидацию, то паникуем
			panic(fmt.Sprintf("config validation error: %v", err))
		}
	}
	return cfg
}
