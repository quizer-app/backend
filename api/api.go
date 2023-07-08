package api

import (
	v1 "github.com/EloToJaa/quizer/api/v1"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1.RegisterRoutes(api)
}
