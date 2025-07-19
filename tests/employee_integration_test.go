package tests

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/common/validator"
	"github.com/nihrom205/idm/inner/database"
	"github.com/nihrom205/idm/inner/employee"
	"github.com/nihrom205/idm/inner/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPageEmployee(t *testing.T) {
	a := assert.New(t)

	cfg := common.GetConfig(".env")
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	employeeRepo := employee.NewEmployeeRepository(db)
	fixture := NewFixtureEmployee(employeeRepo)
	defer func() {
		clearDatabaseEmployee(db)
	}()

	names := []string{"Ivan", "Stas", "Stepan", "Viktor", "Stas2"}
	for _, name := range names {
		fixture.Employee(name)
	}

	app := initTestApp(db)

	type WrappedPageResponse struct {
		Success bool                  `json:"success"`
		Data    employee.PageResponse `json:"data"`
	}

	t.Run("should return 3 employees on first page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=3", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 3)
		a.Equal(int64(5), wrapped.Data.Total)
	})

	t.Run("should return 2 employees on second page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=1&pageSize=3", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 2)
		a.Equal(int64(5), wrapped.Data.Total)
	})

	t.Run("should return 0 employees on third page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=2&pageSize=3", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 0)
		a.Equal(int64(5), wrapped.Data.Total)
	})

	t.Run("should return error for invalid pageSize", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=0", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)

		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		a.Nil(err)
		a.NotNil(errResp)
		a.Contains(errResp["error"], "invalid pageSize")
	})

	t.Run("should use default PageNumber", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageSize=3", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 3)
		a.Equal(int64(5), wrapped.Data.Total)
	})

	t.Run("should use default PageSize", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 1)
		a.Equal(int64(5), wrapped.Data.Total)
	})

	t.Run("should return 2 employee name >= 3 char", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=1&textFilter=sta", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 1)
		a.Equal(int64(2), wrapped.Data.Total)
	})

	t.Run("should return 0 employee name < 3 char", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=3&textFilter=st", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 3)
		a.Equal(int64(5), wrapped.Data.Total)
	})

	t.Run("should return all employee name only spaces", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=3&textFilter=%20%20%20", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.NotNil(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		var wrapped WrappedPageResponse
		err = json.NewDecoder(resp.Body).Decode(&wrapped)
		a.Nil(err)
		a.Len(wrapped.Data.Result, 3)
		a.Equal(int64(5), wrapped.Data.Total)
	})
}

// initTestApp создает Fiber-приложение с реальными зависимостями для интеграционных тестов.
func initTestApp(db *sqlx.DB) *fiber.App {
	cfg := common.GetConfig(".env")

	logger := common.NewLogger(cfg)

	// Валидатор
	vld := validator.NewValidator()

	// Репозиторий и сервис
	employeeRepo := employee.NewEmployeeRepository(db)
	employeeService := employee.NewService(employeeRepo, vld)

	// Создаем сервер и контроллер
	server := web.NewServer()
	employeeController := employee.NewController(server, employeeService, logger)
	employeeController.RegisterRoutes()

	return server.App
}
