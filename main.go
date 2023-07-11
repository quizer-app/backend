package main

import (
	"fmt"
	"log"

	"os"
	"time"

	"github.com/EloToJaa/quizer/api"
	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/initializers"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	initializers.LoadEnvVariables()
	db.ConnectToDb()
}

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Quizer",
		AppName:      "Quizer",
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	app.Static("/", "./public")

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Cookie",
	}))
	// app.Use(csrf.New())
	app.Use(limiter.New(limiter.Config{
		Max:        2000,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests",
			})
		},
	}))

	api.RegisterRoutes(app)

	port := os.Getenv("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

	defer db.DisconnectFromDb()
}
