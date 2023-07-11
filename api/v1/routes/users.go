package routes

import (
	"github.com/EloToJaa/quizer/controllers/users"
	"github.com/EloToJaa/quizer/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router) {
	userRouter := router.Group("/users")

	userRouter.Get("/", middleware.AuthMiddleware, middleware.AdminMiddleware, users.GetUsersController)
}
