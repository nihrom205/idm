package web

import "github.com/gofiber/fiber/v3"

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

	groupInternal := app.Group("/internal")

	groupApi := app.Group("/api")

	// создаём подгруппу "api/v1"
	groupApiV1 := groupApi.Group("/v1")

	return &Server{
		App:           app,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}
