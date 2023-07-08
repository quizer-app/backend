package middleware

import (
	"strings"

	"github.com/EloToJaa/quizer/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	// Get token from header
	token := ctx.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.TrimSpace(token)

	// Check if token is empty
	if token == "" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	// Verify token
	tokenData := &jwt.TokenData{}
	ok, err := tokenData.ParseToken(token, jwt.GetAccessTokenSecret())
	if err != nil || !ok {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	ctx.Locals("data", tokenData)

	return ctx.Next()
}
