package info

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/web"
	"time"
)

type Database interface {
	PingContext(ctx context.Context) error
}

type Controller struct {
	server *web.Server
	cfg    common.Config
	db     Database
}

func NewController(server *web.Server, cfg common.Config, db Database) *Controller {
	return &Controller{
		server: server,
		cfg:    cfg,
		db:     db,
	}
}

type InfoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Controller) RegisterRouters() {
	c.server.GroupInternal.Get("/info", c.GetInfo)
	c.server.GroupInternal.Get("/health", c.GetHealth)
}

// GetInfo получение информации о приложении
func (c *Controller) GetInfo(ctx fiber.Ctx) error {
	resp := &InfoResponse{
		Name:    c.cfg.AppName,
		Version: c.cfg.AppVersion,
	}

	err := ctx.Status(fiber.StatusOK).JSON(resp)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning info")
	}
	return nil
}

// GetHealth проверка работоспособности приложения
func (c *Controller) GetHealth(ctx fiber.Ctx) error {
	// Создаем контекст с таймаутом для проверки БД
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// проверка к подключению БД
	if err := c.db.PingContext(dbCtx); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "Database connection failed")
	}
	return nil
}
