package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ananda1farelas/sewakos-backend/app/model"
	"github.com/ananda1farelas/sewakos-backend/app/service"
	"github.com/ananda1farelas/sewakos-backend/helper"
	"github.com/ananda1farelas/sewakos-backend/middleware"
)

func RegisterAuthRoutes(router fiber.Router, authService service.AuthService) {
	auth := router.Group("/auth")

	auth.Post("/register", func(c *fiber.Ctx) error {
		var req model.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrValidation(fiber.Map{"body": "format request tidak valid"})
		}

		resp, appErr := authService.Register(c.Context(), req)
		if appErr != nil {
			return appErr
		}

		return helper.Success(c, 201, "registrasi berhasil", resp)
	})

	auth.Post("/login", middleware.LoginRateLimiter, func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrValidation(fiber.Map{"body": "format request tidak valid"})
		}

		resp, appErr := authService.Login(c.Context(), req)
		if appErr != nil {
			return appErr
		}

		return helper.Success(c, 200, "login berhasil", resp)
	})

	auth.Post("/refresh", func(c *fiber.Ctx) error {
		var req model.RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrValidation(fiber.Map{"body": "format request tidak valid"})
		}

		resp, appErr := authService.Refresh(c.Context(), req.RefreshToken)
		if appErr != nil {
			return appErr
		}

		return helper.Success(c, 200, "token diperbarui", resp)
	})

	auth.Post("/logout", func(c *fiber.Ctx) error {
		var req model.RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.ErrValidation(fiber.Map{"body": "format request tidak valid"})
		}

		if appErr := authService.Logout(c.Context(), req.RefreshToken); appErr != nil {
			return appErr
		}

		return helper.Success(c, 200, "logout berhasil", nil)
	})

	auth.Get("/me", middleware.RequireAuth, func(c *fiber.Ctx) error {
		userID := helper.GetUserID(c)

		resp, appErr := authService.Me(c.Context(), userID)
		if appErr != nil {
			return appErr
		}

		return helper.Success(c, 200, "data user", resp)
	})
}
