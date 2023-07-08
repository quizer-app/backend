package auth

import (
	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/gofiber/fiber/v2"
)

func LogoutController(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	// Delete refresh token from database mongoDB
	refreshTokenCollection := db.GetCollection(enum.RefreshTokens)
	_, err := refreshTokenCollection.DeleteOne(ctx.Context(), map[string]string{
		"token": refreshToken,
	})
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}
	return nil
}
