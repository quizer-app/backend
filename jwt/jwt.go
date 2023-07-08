package jwt

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(time.Hour * 24 * 7)
}

func GetAccessTokenExpirationTime() time.Time {
	return time.Now().Add(time.Minute * 15)
}

func GetRefreshTokenSecret(passwordHash string) string {
	return os.Getenv("REFRESH_TOKEN_SECRET") + passwordHash
}

func GetAccessTokenSecret() string {
	return os.Getenv("ACCESS_TOKEN_SECRET")
}

func DataFromContext(ctx *fiber.Ctx) *TokenData {
	var data *TokenData = ctx.Locals("data").(*TokenData)
	return data
}
