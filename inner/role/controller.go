package role

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"go.uber.org/zap"
	"slices"
	"strconv"
)

type Controller struct {
	server      *web.Server
	roleService Svc
	logger      *common.Logger
}

// интерфейс сервиса employee.Service
type Svc interface {
	Create(ctx context.Context, request CreateRequest) (int64, error)
	FindById(ctx context.Context, id int64) (Response, error)
	GetAll(ctx context.Context) ([]Response, error)
	FindByIds(ctx context.Context, ids []int64) ([]Response, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByIds(ctx context.Context, ids []int64) error
}

func NewController(server *web.Server, svc Svc, logger *common.Logger) *Controller {
	return &Controller{
		server:      server,
		roleService: svc,
		logger:      logger,
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
// @Description Create a new role.
// @Summary create a new role
// @ID create-role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body role.CreateRequest true "name role"
// @Success 200 {object} common.Response[int64]
// @Failure 400 {object} common.Response[string]
// @Failure 401 {object} common.Response[string]
// @Failure 403 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /role [post]
func (c *Controller) CreateRole(ctx *fiber.Ctx) error {

	// проверяем наличие нужной роли в токене
	claims, err := getClaims(ctx)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}
	if !slices.Contains(claims.RealmAccess.Roles, web.IdmAdmin) {
		return common.ErrResponse(ctx, fiber.StatusForbidden, "Permission denied")
	}

	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("create role: received request", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("create role", zap.Any("request", request))

	// вызываем метод Create сервиса role.Service
	newEmployeeId, err := c.roleService.Create(ctx.Context(), request)
	if err != nil {
		c.logger.Error("create role", zap.Any("request", request))
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
		c.logger.Error("create role", zap.Any("request", request))
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
		return err
	}
	return nil
}

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/roles/:id"
// @Description Get role.
// @Summary get role
// @ID get-role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "id role"
// @Success 200 {object} common.Response[role.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 401 {object} common.Response[string]
// @Failure 403 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /roles/{id} [get]
func (c *Controller) GetRole(ctx *fiber.Ctx) error {

	// проверяем наличие нужной роли в токене
	claims, err := getClaims(ctx)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}
	if !slices.Contains(claims.RealmAccess.Roles, web.IdmAdmin) &&
		!slices.Contains(claims.RealmAccess.Roles, web.IdmUser) {
		return common.ErrResponse(ctx, fiber.StatusForbidden, "Permission denied")
	}

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	c.logger.Debug("get role", zap.Any("idParam", idParam))
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.logger.Error("get role: invalid id param", zap.Any("idParam", idParam))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод FindById сервиса role.Service
	response, err := c.roleService.FindById(ctx.Context(), id)
	if err != nil {
		c.logger.Error("get role", zap.Any("request", idParam))
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
		c.logger.Error("get role", zap.Any("request", idParam))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/roles"
// @Description Get all role.
// @Summary get all role
// @ID get-all-role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response[role.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 401 {object} common.Response[string]
// @Failure 403 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /roles [get]
func (c *Controller) GetAllRoles(ctx *fiber.Ctx) error {

	// проверяем наличие нужной роли в токене
	claims, err := getClaims(ctx)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}
	if !slices.Contains(claims.RealmAccess.Roles, web.IdmAdmin) &&
		!slices.Contains(claims.RealmAccess.Roles, web.IdmUser) {
		return common.ErrResponse(ctx, fiber.StatusForbidden, "Permission denied")
	}

	// вызываем метод GetAll сервиса role.Service
	response, err := c.roleService.GetAll(ctx.Context())
	if err != nil {
		c.logger.Error("get all roles", zap.Any("request", err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		c.logger.Error("get all roles", zap.Any("request", err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/roles/ids"
// @Description Get role by id.
// @Summary get role by id
// @ID get-role-by-id
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ids body role.FindByIdsRequest true "ids role"
// @Success 200 {object} common.Response[role.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 401 {object} common.Response[string]
// @Failure 403 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /roles/ids [post]
func (c *Controller) GetRoleByIds(ctx *fiber.Ctx) error {

	// проверяем наличие нужной роли в токене
	claims, err := getClaims(ctx)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}
	if !slices.Contains(claims.RealmAccess.Roles, web.IdmAdmin) &&
		!slices.Contains(claims.RealmAccess.Roles, web.IdmUser) {
		return common.ErrResponse(ctx, fiber.StatusForbidden, "Permission denied")
	}

	// анмаршалим JSON body запроса в структуру FindByIdsRequest
	var request FindByIdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("get role: received request", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("get role by ids", zap.Any("request", request))

	// вызываем метод FindByIds сервиса role.Service
	response, err := c.roleService.FindByIds(ctx.Context(), request.Ids)
	if err != nil {
		c.logger.Error("get role", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, response); err != nil {
		c.logger.Error("get role", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/roles/:id"
// @Description Delete role by id.
// @Summary delete role by id
// @ID delete-role-by-id
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "id role"
// @Success 200 {object} common.Response[int64]
// @Failure 400 {object} common.Response[string]
// @Failure 401 {object} common.Response[string]
// @Failure 403 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /role/{id} [delete]
func (c *Controller) DeleteRole(ctx *fiber.Ctx) error {

	// проверяем наличие нужной роли в токене
	claims, err := getClaims(ctx)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}
	if !slices.Contains(claims.RealmAccess.Roles, web.IdmAdmin) {
		return common.ErrResponse(ctx, fiber.StatusForbidden, "Permission denied")
	}

	// получаем ID из параметра маршрута
	idParam := ctx.Params("id")
	c.logger.Debug("delete role", zap.Any("idParam", idParam))
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.logger.Error("delete role: invalid id param", zap.Any("idParam", idParam))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid employee id")
	}

	// вызываем метод DeleteById сервиса role.Service
	err = c.roleService.DeleteById(ctx.Context(), id)
	if err != nil {
		c.logger.Error("delete role", zap.Any("request", idParam))
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
		c.logger.Error("delete role", zap.Any("request", idParam))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/roles/ids"
// @Description Delete role by list ids.
// @Summary delete role by list ids
// @ID delete-role-by-list-ids
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ids body role.DeleteByIdsRequest true "ids role"
// @Success 200 {object} common.Response[int64]
// @Failure 400 {object} common.Response[string]
// @Failure 401 {object} common.Response[string]
// @Failure 403 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /role/ids [delete]
func (c *Controller) DeleteRolesByIds(ctx *fiber.Ctx) error {

	// проверяем наличие нужной роли в токене
	claims, err := getClaims(ctx)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}
	if !slices.Contains(claims.RealmAccess.Roles, web.IdmAdmin) {
		return common.ErrResponse(ctx, fiber.StatusForbidden, "Permission denied")
	}

	// анмаршалим JSON body запроса в структуру DeleteByIdsRequest
	var request DeleteByIdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("delete roles: received request", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("delete roles by ids", zap.Any("request", request))

	// вызываем метод DeleteByIds сервиса role.Service
	err = c.roleService.DeleteByIds(ctx.Context(), request.Ids)
	if err != nil {
		c.logger.Error("delete roles", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// возвращаем успешный ответ
	if err := common.OkResponse(ctx, struct{}{}); err != nil {
		c.logger.Error("delete roles", zap.Any("request", request))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func getClaims(ctx *fiber.Ctx) (*web.IdmClaims, error) {
	token, ok := ctx.Locals(web.JwtKey).(*jwt.Token)
	if !ok || token == nil {
		return nil, errors.New("missing or invalid token")
	}
	claims, ok := token.Claims.(*web.IdmClaims)
	if !ok || claims == nil {
		return nil, errors.New("missing or invalid claims")
	}
	return claims, nil
}
