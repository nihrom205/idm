package web

import "github.com/gofiber/fiber/v2"

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

	groupInternal := app.Group("/internal")

	groupApi := app.Group("/api")

	groupApiV1 := groupApi.Group("/v1")

	return &Server{
		App:           app,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}
