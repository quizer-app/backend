package middleware

import (
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/jwt"
	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(ctx *fiber.Ctx) error {
	data := jwt.DataFromContext(ctx)

	if data.User.Role != enum.Admin {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return ctx.Next()
}
