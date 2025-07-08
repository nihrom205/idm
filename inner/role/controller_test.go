package role

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Объявляем структуру мока сервиса employee.Service
type MockService struct {
	mock.Mock
}

// Реализуем функции мок-сервиса
func (svc *MockService) FindById(ctx context.Context, id int64) (Response, error) {
	args := svc.Called(id)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) Create(ctx context.Context, request CreateRequest) (int64, error) {
	args := svc.Called(request)
	return args.Get(0).(int64), args.Error(1)
}

func (svc *MockService) GetAll(ctx context.Context) ([]Response, error) {
	args := svc.Called()
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) FindByIds(ctx context.Context, ids []int64) ([]Response, error) {
	args := svc.Called(ids)
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) DeleteById(ctx context.Context, id int64) error {
	args := svc.Called(id)
	return args.Error(0)
}

func (svc *MockService) DeleteByIds(ctx context.Context, ids []int64) error {
	args := svc.Called(ids)
	return args.Error(0)
}

func TestController_CreateRole(t *testing.T) {
	var a = assert.New(t)
	// Создаем тестовый логгер
	logger := &common.Logger{
		Logger: zap.NewNop(), // Логгер, который ничего не делает (подходит для тестов),
	}

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should return created role id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"name\": \"john doe\"}")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("Create", mock.AnythingOfType("CreateRequest")).Return(int64(123), nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(123), responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should return error if role already exists", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"name\": \"john doe\"}")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("Create", mock.AnythingOfType("CreateRequest")).
			Return(int64(0), common.AlreadyExistsError{Message: "employee already exists"})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(0), responseBody.Data)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return error validator error", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"name\": \"john doe\"}")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("Create", mock.AnythingOfType("CreateRequest")).
			Return(int64(0), common.RequestValidatorError{Message: "employee validation error"})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(0), responseBody.Data)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return error transaction", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"name\": \"john doe\"}")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("Create", mock.AnythingOfType("CreateRequest")).Return(int64(0), errors.New("transaction error"))

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(0), responseBody.Data)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return error unmarshal", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("invalid json")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("Create", mock.AnythingOfType("CreateRequest")).Return(int64(0), errors.New("marshal error"))

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(0), responseBody.Data)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})
}

func TestController_GetEmployee(t *testing.T) {
	var a = assert.New(t)
	// Создаем тестовый логгер
	logger := &common.Logger{
		Logger: zap.NewNop(), // Логгер, который ничего не делает (подходит для тестов),
	}

	t.Run("should return employee", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		req.Header.Set("Content-Type", "application/json")

		response := Response{
			Id:       123,
			Name:     "john doe",
			CreateAt: time.Time{},
			UpdateAt: time.Time{},
		}

		// Настраиваем поведение мока в тесте
		svc.On("FindById", int64(123)).Return(response, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(response, responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should return err bad id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/bad", nil)
		req.Header.Set("Content-Type", "application/json")

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return err validation", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("FindById", int64(123)).Return(Response{}, common.RequestValidatorError{Message: "invalid employee"})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return err transaction", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("FindById", int64(123)).Return(Response{}, common.RepositoryError{Message: "transaction error"})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return err not found", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("FindById", int64(123)).Return(Response{}, common.NotFoundError{Message: "employee not found"})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusNotFound, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should return err", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("FindById", int64(123)).Return(Response{}, errors.New("error"))

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})
}

func TestController_GetAllRoles(t *testing.T) {
	var a = assert.New(t)
	// Создаем тестовый логгер
	logger := &common.Logger{
		Logger: zap.NewNop(), // Логгер, который ничего не делает (подходит для тестов),
	}

	t.Run("should success get all roles", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/", nil)
		req.Header.Set("Content-Type", "application/json")

		responses := []Response{
			Response{
				Id:       123,
				Name:     "john doe",
				CreateAt: time.Time{},
				UpdateAt: time.Time{},
			},
			Response{
				Id:       124,
				Name:     "dred bev",
				CreateAt: time.Time{},
				UpdateAt: time.Time{},
			},
		}

		// Настраиваем поведение мока в тесте
		svc.On("GetAll").Return(responses, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(responses, responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should return err", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/", nil)
		req.Header.Set("Content-Type", "application/json")

		responses := []Response{}

		// Настраиваем поведение мока в тесте
		svc.On("GetAll").Return(responses, errors.New("error"))

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Nil(responseBody.Data)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})
}

func TestController_GetEmployeeByIds(t *testing.T) {
	var a = assert.New(t)
	// Создаем тестовый логгер
	logger := &common.Logger{
		Logger: zap.NewNop(), // Логгер, который ничего не делает (подходит для тестов),
	}

	t.Run("should success get all roles by ids", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"ids\": [123,124]}")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles/ids", body)
		req.Header.Set("Content-Type", "application/json")

		responses := []Response{
			Response{
				Id:       123,
				Name:     "john doe",
				CreateAt: time.Time{},
				UpdateAt: time.Time{},
			},
			Response{
				Id:       124,
				Name:     "dred bev",
				CreateAt: time.Time{},
				UpdateAt: time.Time{},
			},
		}

		// Настраиваем поведение мока в тесте
		svc.On("FindByIds", []int64{123, 124}).Return(responses, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(len(responses), len(responseBody.Data))
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should return err", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"ids\": [bad, bad2]}")
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/roles/ids", body)
		req.Header.Set("Content-Type", "application/json")

		responses := []Response{}

		// Настраиваем поведение мока в тесте
		svc.On("FindByIds", []int64{123, 124}).Return(responses, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(len(responses), len(responseBody.Data))
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})
}

func TestController_DeleteEmployee(t *testing.T) {
	var a = assert.New(t)
	// Создаем тестовый логгер
	logger := &common.Logger{
		Logger: zap.NewNop(), // Логгер, который ничего не делает (подходит для тестов),
	}

	t.Run("should success del by id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/123", nil)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("DeleteById", int64(123)).Return(nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should success err invalid id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/badId", nil)
		req.Header.Set("Content-Type", "application/json")

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})
}

func TestController_DeleteRolesByIds(t *testing.T) {
	var a = assert.New(t)
	// Создаем тестовый логгер
	logger := &common.Logger{
		Logger: zap.NewNop(), // Логгер, который ничего не делает (подходит для тестов),
	}

	t.Run("should success del by ids", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"ids\": [123,124]}")
		req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/ids", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("DeleteByIds", []int64{123, 124}).Return(nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should err bad request", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		body := strings.NewReader("{\"ids\": [bad,fail]}")
		req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/ids", body)
		req.Header.Set("Content-Type", "application/json")

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})

	t.Run("should success del by ids", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		svc := &MockService{}
		controller := NewController(server, svc, logger)
		controller.RegisterRoutes()

		// Готовим тестовое окружение
		//request := DeleteByIdsRequest{
		//	Ids: []int64{123, 124},
		//}
		body := strings.NewReader("{\"ids\": [123,124]}")
		req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/ids", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("DeleteByIds", []int64{123, 124}).Return(errors.New("error transaction"))

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
	})
}
