package employee

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"go.uber.org/zap"
	"strconv"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
	logger          *common.Logger
}

// интерфейс сервиса employee.Service
type Svc interface {
	Create(ctx context.Context, request CreateRequest) (int64, error)
	FindById(ctx context.Context, id int64) (Response, error)
	GetAll(ctx context.Context) ([]Response, error)
	FindByIds(ctx context.Context, ids []int64) ([]Response, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByIds(ctx context.Context, ids []int64) error
	FindPage(ctx context.Context, req PageRequest) (PageResponse, error)
}

func NewController(server *web.Server, svc Svc, logger *common.Logger) *Controller {
	return &Controller{
		server:          server,
		employeeService: svc,
		logger:          logger,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/page", c.GetPageEmployee)
	c.server.GroupApiV1.Get("/employees/:id", c.GetEmployee)
	c.server.GroupApiV1.Get("/employees", c.GetAllEmployees)
	c.server.GroupApiV1.Post("/employees/ids", c.GetEmployeeByIds)
	c.server.GroupApiV1.Delete("/employees/ids", c.DeleteEmployeesByIds)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteEmployee)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (c *Controller) CreateEmployee(ctx *fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// логируем тело запроса
	c.logger.Debug("create employee: received request", zap.Any("request", request))
	// вызываем метод Create сервиса employee.Service
	newEmployeeId, err := c.employeeService.Create(ctx.Context(), request)
	if err != nil {
		c.logger.Error("create employee", zap.Error(err))
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
		c.logger.Error("create employee", zap.Error(err))
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
		return err
	}
	return nil
}

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees/:id"
func (c *Controller) GetEmployee(ctx *fiber.Ctx) error {

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	c.logger.Debug("get employee", zap.String("id", idParam))
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.logger.Error("get employee", zap.String("id", idParam), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод FindById сервиса employee.Service
	response, err := c.employeeService.FindById(ctx.Context(), id)
	if err != nil {
		c.logger.Error("get employee", zap.String("id", idParam), zap.Error(err))
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
		c.logger.Error("get employee", zap.String("id", idParam), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees"
func (c *Controller) GetAllEmployees(ctx *fiber.Ctx) error {

	// вызываем метод GetAll сервиса employee.Service
	response, err := c.employeeService.GetAll(ctx.Context())
	if err != nil {
		c.logger.Error("get all employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		c.logger.Error("get all employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees/ids"
func (c *Controller) GetEmployeeByIds(ctx *fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру FindByIdsRequest
	var request FindByIdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("get employee by ids", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("get employee by ids", zap.Any("request", request))

	// вызываем метод FindByIds сервиса employee.Service
	response, err := c.employeeService.FindByIds(ctx.Context(), request.Ids)
	if err != nil {
		c.logger.Error("get employee by ids", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		c.logger.Error("get employee by ids", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/employees/:id"
func (c *Controller) DeleteEmployee(ctx *fiber.Ctx) error {

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	c.logger.Debug("delete employee", zap.String("id", idParam))
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.logger.Error("delete employee", zap.String("id", idParam), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод DeleteById сервиса employee.Service
	err = c.employeeService.DeleteById(ctx.Context(), id)
	if err != nil {
		c.logger.Error("delete employee", zap.String("id", idParam), zap.Error(err))
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
		c.logger.Error("delete employee", zap.String("id", idParam), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/employees/ids"
func (c *Controller) DeleteEmployeesByIds(ctx *fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру DeleteByIdsRequest
	var request DeleteByIdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("delete employees by ids", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("delete employees by ids", zap.Any("request", request))

	// вызываем метод DeleteByIds сервиса employee.Service
	err := c.employeeService.DeleteByIds(ctx.Context(), request.Ids)
	if err != nil {
		c.logger.Error("delete employees by ids", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, struct{}{}); err != nil {
		c.logger.Error("delete employees by ids", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// GetEmployeesPage получает страницу сотрудников
// функция-хендлер, которая будет вызываться при GET запросе по маршруту /api/v1/employees/page?pageNumber=x&pageSize=y
func (c *Controller) GetPageEmployee(ctx *fiber.Ctx) error {
	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "1"))
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid pageSize")
	}

	pageNumber, err := strconv.Atoi(ctx.Query("pageNumber", "0"))
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	textFilter := ctx.Query("textFilter", "")

	// идем в бд за данными
	request := PageRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		TextFilter: textFilter,
	}
	c.logger.Debug("get page employee by pageNumber and pageSize", zap.Any("request", request))

	page, err := c.employeeService.FindPage(ctx.Context(), request)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid pageSize")
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, page); err != nil {
		c.logger.Error("get page employee by pageNumber and pageSize", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
