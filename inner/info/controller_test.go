package info

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http/httptest"
	"testing"
)

// MockDatabase - мок для интерфейса Database
type MockDatabase struct {
	mock.Mock
}

// Реализует интерфейс Database
func (m *MockDatabase) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// setupTest инициализирует тестовое окружение
func setupTest(t *testing.T) (*fiber.App, *MockDatabase) {
	// Готовим тестовое окружение
	server := web.NewServer()
	cfg := common.Config{
		DbDriverName: "postgres",
		DSN:          "test_dsn",
		AppName:      "test_app",
		AppVersion:   "0.0.1",
	}
	mock := &MockDatabase{}
	controller := NewController(server, cfg, mock)

	if controller == nil {
		t.Fatal("Failed to create controller")
	}

	if server.GroupInternal == nil {
		t.Fatal("GroupInternal is nil")
	}

	controller.RegisterRouters()
	return server.App, mock
}

func TestName(t *testing.T) {
	var a = assert.New(t)

	t.Run("Success", func(t *testing.T) {
		app, _ := setupTest(t)
		req := httptest.NewRequest(fiber.MethodGet, "/internal/info", nil)
		resp, err := app.Test(req)
		a.Nil(err)

		a.Equal(fiber.StatusOK, resp.StatusCode)

		var responseBody InfoResponse
		bytesData, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("Failed to read response body")
		}
		err = json.Unmarshal(bytesData, &responseBody)
		if err != nil {
			t.Fatal("Failed to unmarshal response body")
		}

		a.Equal("test_app", responseBody.Name)
		a.Equal("0.0.1", responseBody.Version)
	})
}
func TestGetHealth(t *testing.T) {
	var a = assert.New(t)

	t.Run("Success - Database available", func(t *testing.T) {
		app, mockDB := setupTest(t)
		mockDB.On("PingContext", mock.Anything).Return(nil)

		req := httptest.NewRequest(fiber.MethodGet, "/internal/health", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.Equal(fiber.StatusOK, resp.StatusCode)

		mockDB.AssertCalled(t, "PingContext", mock.Anything)
	})

	t.Run("Success - Database unavailable", func(t *testing.T) {
		app, mockDB := setupTest(t)
		mockDB.On("PingContext", mock.Anything).Return(errors.New("database connection failed"))

		req := httptest.NewRequest(fiber.MethodGet, "/internal/health", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.Equal(fiber.StatusInternalServerError, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		var errResponse common.Response[string]
		if err != nil {
			t.Fatal("Failed to read response body")
		}
		err = json.Unmarshal(bytesData, &errResponse)
		if err != nil {
			t.Fatal("Failed to unmarshal response body")
		}
		a.False(errResponse.Success)
		a.Equal("Database connection failed", errResponse.Message)

		mockDB.AssertCalled(t, "PingContext", mock.Anything)
	})
}
