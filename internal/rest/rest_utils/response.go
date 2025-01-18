package rest_utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhrifqii/tuskar/domain"
)

func ReturnOkResponse[T interface{}](c *fiber.Ctx, data T) error {
	r := domain.ApiResponse[T]{
		Message: "Success",
		Data:    data,
	}
	return c.Status(fiber.StatusOK).JSON(r)
}

func ReturnAnyResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}
