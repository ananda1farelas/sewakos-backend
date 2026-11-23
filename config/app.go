package config

import (
	"github.com/ananda1farelas/sewakos-backend/app/model"
	"github.com/ananda1farelas/sewakos-backend/helper"
	"github.com/gofiber/fiber/v2"
)

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: centralErrorHandler,
	})
	return app
}

func centralErrorHandler(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*helper.AppError); ok {
		return c.Status(appErr.Status).JSON(model.WebResponse{
			Success: false,
			Message: appErr.Message,
			Errors: fiber.Map{
				"code":    appErr.Code,
				"details": appErr.Details,
			},
		})
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(model.WebResponse{
			Success: false,
			Message: fiberErr.Message,
			Errors: fiber.Map{
				"code": helper.CodeInternal,
			},
		})
	}

	return c.Status(500).JSON(model.WebResponse{
		Success: false,
		Message: "terjadi kesalahan pada server",
		Errors: fiber.Map{
			"code": helper.CodeInternal,
		},
	})
}
