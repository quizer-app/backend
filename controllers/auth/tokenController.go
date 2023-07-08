package auth

import (
	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/jwt"
	"github.com/EloToJaa/quizer/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func TokenController(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	refreshTokenModel := &models.RefreshToken{}

	// Check if refresh token exists in database
	refreshTokenCollection := db.GetCollection(enum.RefreshTokens)
	err := refreshTokenCollection.FindOne(ctx.Context(), bson.M{"token": refreshToken}).Decode(&refreshTokenModel)
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	// Verify refresh token
	accessTokenData := &jwt.TokenData{}
	valid, err := accessTokenData.ParseToken(refreshToken, jwt.GetRefreshTokenSecret(refreshTokenModel.UserPassword))
	if err != nil || !valid {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	// Generate new access token
	accessToken, err := accessTokenData.GenerateToken(jwt.GetAccessTokenSecret())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
		"message":      "Token refreshed",
	})

	return nil
}
