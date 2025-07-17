package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/nihrom205/idm/docs" // обязательно импортируем наш пакет с документацией
)

type Server struct {
	App *fiber.App
	// группа публичного API
	GroupApiV1 fiber.Router
	// группа непубличного API
	GroupInternal fiber.Router
}

func NewServer() *Server {
	// создаём новый веб-вервер
	app := fiber.New()

	// подключаем middleware
	registerMiddleware(app)

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	groupInternal := app.Group("/internal")

	groupApi := app.Group("/api")

	groupApiV1 := groupApi.Group("/v1")

	return &Server{
		App:           app,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}
