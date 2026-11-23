package helper

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type RequestContext struct {
	Ctx    context.Context
	UserID int
	Role   string
}

func NewRequestContext(c *fiber.Ctx) *RequestContext {
	userID, _ := c.Locals("user_id").(int)
	role, _ := c.Locals("role").(string)

	return &RequestContext{
		Ctx:    c.Context(),
		UserID: userID,
		Role:   role,
	}
}

func ParamID(c *fiber.Ctx, key string) (int, error) {
	id, err := c.ParamsInt(key)
	if err != nil {
		return 0, err
	}
	return id, nil
}
