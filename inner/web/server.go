package web

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/nihrom205/idm/docs" // обязательно импортируем наш пакет с документацией
)

type Server struct {
	App *fiber.App
	// группа публичного API
	GroupApi fiber.Router
	// группа публичного API первой версии
	GroupApiV1 fiber.Router
	// группа непубличного API
	GroupInternal fiber.Router
}

type AuthMiddlewareInterface interface {
	ProtectWithJwt() func(*fiber.Ctx) error
}

func NewServer() *Server {
	// создаём новый веб-вервер
	app := fiber.New()

	// подключаем middleware
	registerMiddleware(app)

	groupInternal := app.Group("/internal")

	groupApi := app.Group("/api")

	groupApiV1 := groupApi.Group("/v1")

	return &Server{
		App:           app,
		GroupApi:      groupApi,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}
