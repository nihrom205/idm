package role

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"strconv"
)

type Controller struct {
	server      *web.Server
	roleService Svc
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
		server:      server,
		roleService: svc,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/roles", c.CreateRole)
	c.server.GroupApiV1.Get("/roles/:id", c.GetRole)
	c.server.GroupApiV1.Get("/roles", c.GetAllRoles)
	c.server.GroupApiV1.Post("/roles/ids", c.GetRoleByIds)
	c.server.GroupApiV1.Delete("/roles/ids", c.DeleteRolesByIds)
	c.server.GroupApiV1.Delete("/roles/:id", c.DeleteRole)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/role"
func (c *Controller) CreateRole(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод Create сервиса role.Service
	newEmployeeId, err := c.roleService.Create(request)
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

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/roles/:id"
func (c *Controller) GetRole(ctx fiber.Ctx) error {

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод FindById сервиса role.Service
	response, err := c.roleService.FindById(id)
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

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/roles"
func (c *Controller) GetAllRoles(ctx fiber.Ctx) error {

	// вызываем метод GetAll сервиса role.Service
	response, err := c.roleService.GetAll()
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/roles/ids"
func (c *Controller) GetRoleByIds(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру FindByIdsRequest
	var request FindByIdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод FindByIds сервиса role.Service
	response, err := c.roleService.FindByIds(request.Ids)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/roles/:id"
func (c *Controller) DeleteRole(ctx fiber.Ctx) error {

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод DeleteById сервиса role.Service
	err = c.roleService.DeleteById(id)
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

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/roles/ids"
func (c *Controller) DeleteRolesByIds(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру DeleteByIdsRequest
	var request DeleteByIdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод DeleteByIds сервиса role.Service
	err := c.roleService.DeleteByIds(request.Ids)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, struct{}{}); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
