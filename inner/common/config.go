package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName   string `validate:"required"`
	DSN            string `validate:"required"`
	AppName        string `validate:"required"`
	AppVersion     string `validate:"required"`
	LogLevel       string `validate:"required"`
	LogDevelopMode bool   `validate:"required"`
	SslCert        string `validate:"required"`
	SslKey         string `validate:"required"`
	KeycloakJwkUrl string `validate:"required"`
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) Config {
	// если нет файла, то залогируем это и попробуем получить конфиг из переменных окружения
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", zap.Error(err))
	}

	cfg := Config{
		DbDriverName:   os.Getenv("DB_DRIVER_NAME"),
		DSN:            os.Getenv("DB_DSN"),
		AppName:        os.Getenv("APP_NAME"),
		AppVersion:     os.Getenv("APP_VERSION"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
		LogDevelopMode: os.Getenv("LOG_DEVELOP_MODE") == "true",
		SslCert:        os.Getenv("SSL_CERT"),
		SslKey:         os.Getenv("SSL_KEY"),
		KeycloakJwkUrl: os.Getenv("KEYCLOAK_JWK_URL"),
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
