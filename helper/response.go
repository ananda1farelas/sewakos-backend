package helper

import (
	"github.com/ananda1farelas/sewakos-backend/app/model"
	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *fiber.Ctx, status int, message string, data interface{}, meta *model.Meta) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Fail(c *fiber.Ctx, status int, message string, errors interface{}) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
