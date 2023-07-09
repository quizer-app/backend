package routes

import (
	"github.com/EloToJaa/quizer/controllers/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	authRouter := router.Group("/auth")

	authRouter.Post("/login", auth.LoginController)
	authRouter.Post("/register", auth.RegisterController)
	authRouter.Post("/token", auth.TokenController)
	authRouter.Delete("/logout", auth.LogoutController)
	authRouter.Post("/verify", auth.VerifyController)
	authRouter.Post("/verify/:id", auth.VerifyUserController)
	authRouter.Post("/forgot-password", auth.ForgotPasswordController)
	authRouter.Post("/reset-password/:id", auth.ResetPasswordController)
}
