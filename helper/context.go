package helper

import "github.com/gofiber/fiber/v2"

func SetAuthLocals(c *fiber.Ctx, userID int, role string) {
	c.Locals("user_id", userID)
	c.Locals("role", role)
}

func GetUserID(c *fiber.Ctx) int {
	id, _ := c.Locals("user_id").(int)
	return id
}

func GetRole(c *fiber.Ctx) string {
	role, _ := c.Locals("role").(string)
	return role
}
