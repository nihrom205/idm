package web

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	a := assert.New(t)

	t.Run("should recover from panic and return 500", func(t *testing.T) {
		server := NewServer()

		// Настраиваем роут для теста
		route := func(app *fiber.App) {
			app.Get("/panic", func(c fiber.Ctx) error {
				panic("test panic")
			})
		}
		route(server.App)

		// Создаем тестовый запрос
		req := httptest.NewRequest("GET", "/panic", nil)
		resp, err := server.App.Test(req)

		// Проверяем результат
		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		// Читаем тело ответа
		body, err := io.ReadAll(resp.Body)
		a.Nil(err)
		a.Contains(string(body), "test panic")
	})
}

func TestRequestIdMiddleware(t *testing.T) {
	a := assert.New(t)

	t.Run("should use provided request ID", func(t *testing.T) {
		server := NewServer()

		// Добавляем тестовый роут, который возвращает Request ID
		server.App.Get("/test", func(c fiber.Ctx) error {
			requestId := c.Locals("requestid")
			return c.JSON(fiber.Map{
				"request_id": requestId,
			})
		})

		// Создаем запрос
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-request-id-0000")

		// Выполняем запрос
		resp, err := server.App.Test(req)
		a.Nil(err)

		// Проверяем статус
		a.Equal(http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		a.Nil(err)
		a.NotNil(body)
		a.NotNil(response["request_id"])
		a.Equal("test-request-id-0000", response["request_id"])

		// Проверяем, что заголовок X-Request-ID установлен в ответе
		responseRequestId := resp.Header.Get("X-Request-ID")
		a.Equal("test-request-id-0000", responseRequestId)
	})

	t.Run("should use note provided request ID", func(t *testing.T) {
		server := NewServer()

		// Добавляем тестовый роут, который возвращает Request ID
		server.App.Get("/test", func(c fiber.Ctx) error {
			requestId := c.Locals("requestid")
			return c.JSON(fiber.Map{
				"request_id": requestId,
			})
		})

		// Создаем запрос
		req := httptest.NewRequest("GET", "/test", nil)

		// Выполняем запрос
		resp, err := server.App.Test(req)
		a.Nil(err)

		// Проверяем статус
		a.Equal(http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		a.Nil(err)
		a.NotNil(body)
		a.NotNil(response["request_id"])

		//Проверяем, что заголовок X-Request-ID установлен в ответе
		responseRequestId := resp.Header.Get("X-Request-ID")
		a.NotNil(responseRequestId)
	})
}
