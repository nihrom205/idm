package employee

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"strconv"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
}

// интерфейс сервиса employee.Service
type Svc interface {
	Create(request CreateRequest) (int64, error)
	FindById(id int64) (Response, error)
	GetAll() ([]Response, error)
	FindByIds(ids []int64) ([]Response, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

func NewController(server *web.Server, svc Svc) *Controller {
	return &Controller{
		server:          server,
		employeeService: svc,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/:id", c.GetEmployee)
	c.server.GroupApiV1.Get("/employees", c.GetAllEmployees)
	c.server.GroupApiV1.Post("/employees/ids", c.GetEmployeeByIds)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteEmployee)
	c.server.GroupApiV1.Delete("/employees/ids", c.DeleteEmployeesByIds)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (c *Controller) CreateEmployee(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод CreateEmployee сервиса employee.Service
	newEmployeeId, err := c.employeeService.Create(request)
	if err != nil {
		switch {
		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidatorError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newEmployeeId); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
		return err
	}
	return nil
}

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees/:id"
func (c *Controller) GetEmployee(ctx fiber.Ctx) error {

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод FindById сервиса employee.Service
	response, err := c.employeeService.FindById(id)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidatorError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.RepositoryError{}):
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrResponse(ctx, fiber.StatusNotFound, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees"
func (c *Controller) GetAllEmployees(ctx fiber.Ctx) error {

	// вызываем метод GetAll сервиса employee.Service
	response, err := c.employeeService.GetAll()
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees/ids"
func (c *Controller) GetEmployeeByIds(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру FindEmployeesByIdsRequest
	var request FindEmployeesByIdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод FindByIds сервиса employee.Service
	response, err := c.employeeService.FindByIds(request.Ids)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/employees/:id"
func (c *Controller) DeleteEmployee(ctx fiber.Ctx) error {
	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод FindById сервиса employee.Service
	err = c.employeeService.DeleteById(id)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidatorError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.RepositoryError{}):
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrResponse(ctx, fiber.StatusNotFound, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	if err := common.OkResponse(ctx, struct{}{}); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/employees/ids"
func (c *Controller) DeleteEmployeesByIds(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру DeleteEmployeesByIdsRequest
	var request DeleteEmployeesByIdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод DeleteByIds сервиса employee.Service
	err := c.employeeService.DeleteByIds(request.Ids)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, struct{}{}); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
