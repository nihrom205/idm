package common

import (
	"github.com/gofiber/fiber/v2"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"error"`
	Data    T      `json:"data"`
}

func ErrResponse(c *fiber.Ctx, code int, message string) error {
	res := Response[any]{
		Success: false,
		Message: message,
		Data:    nil,
	}
	return c.Status(code).JSON(res)
}

func OkResponse[T any](c *fiber.Ctx, data T) error {
	return c.JSON(&Response[T]{
		Success: true,
		Data:    data,
	})
}
