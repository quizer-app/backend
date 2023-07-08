package auth

import (
	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/models"
	"github.com/EloToJaa/quizer/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

type ForgotPasswordForm struct {
	Email string `json:"email" validate:"required,email,min=3,max=64"`
}

func ForgotPasswordController(ctx *fiber.Ctx) error {
	var body ForgotPasswordForm

	// Parse body into struct
	if err := ctx.BodyParser(&body); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
		})
		return err
	}

	// Validate body
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  utils.FormatValidationErrors(err),
		})
		return nil
	}

	userCollection := db.GetCollection(enum.Users)
	userModel := &models.User{}

	// Check if user exists
	err := userCollection.FindOne(ctx.Context(), bson.M{"email": body.Email}).Decode(&userModel)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "User not found",
		})
		return nil
	}

	go utils.ResetPassword(userModel)
	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email sent",
	})

	return nil
}
