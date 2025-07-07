package web

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/uuid"
)

func registerMiddleware(app *fiber.App) {
	app.Use(recover.New())
	//app.Use(requestid.New())
	app.Use(func(c fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		// Сохраняем в locals, чтобы потом можно было получить из handler
		c.Locals("requestid", requestID)
		// Добавляем в заголовок ответа
		c.Set("X-Request-ID", requestID)
		return c.Next()
	})
	app.Use(logger.New())
}
