package middleware

import (
	"log/slog"
	"time"

	"github.com/ananda1farelas/sewakos-backend/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Register(app *fiber.App) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(RequestLogger())
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		slog.Info("request",
			"request_id", c.Locals("requestid"),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
		)

		return err
	}
}

func RequireJSON(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet || c.Method() == fiber.MethodDelete {
		return c.Next()
	}

	if c.Get("Content-Type") != fiber.MIMEApplicationJSON {
		return helper.ErrValidation("Content-Type harus application/json")
	}

	return c.Next()
}
