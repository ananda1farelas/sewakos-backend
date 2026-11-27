package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ananda1farelas/sewakos-backend/database"
	"github.com/ananda1farelas/sewakos-backend/helper"
)

func RegisterHealth(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		ctx := c.Context()

		if err := database.DB.Ping(ctx); err != nil {
			return helper.Fail(c, 503, "database tidak sehat", fiber.Map{"postgres": "down"})
		}

		return helper.Success(c, 200, "server berjalan", fiber.Map{
			"postgres": "up",
		})
	})
}

func APIGroup(app *fiber.App) fiber.Router {
	return app.Group("/api/v1")
}
