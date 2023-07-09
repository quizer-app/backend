package v1

import (
	"github.com/EloToJaa/quizer/api/v1/routes"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")

	routes.AuthRoutes(v1)
	routes.UserRoutes(v1)
}
