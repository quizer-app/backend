package users

import (
	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsersController(ctx *fiber.Ctx) error {
	// Get all users from database
	userCollection := db.GetCollection(enum.Users)

	// Get all users
	var users []models.User
	cursor, err := userCollection.Find(ctx.Context(), bson.M{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if err := cursor.All(ctx.Context(), &users); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Return users
	return ctx.JSON(fiber.Map{
		"message": "Success",
		"users":   users,
	})
}
