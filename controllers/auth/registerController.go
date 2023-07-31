package auth

import (
	"time"

	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/models"
	"github.com/EloToJaa/quizer/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

type RegisterForm struct {
	Username        string `json:"username" validate:"required,min=3,max=64"`
	Email           string `json:"email" validate:"required,email,min=3,max=64"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=64"`
}

func RegisterController(ctx *fiber.Ctx) error {
	var body RegisterForm

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
	err := userCollection.FindOne(ctx.Context(), bson.M{"username": body.Username}).Decode(&userModel)
	if err == nil {
		ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
		return nil
	}

	// Check if email exists
	err = userCollection.FindOne(ctx.Context(), bson.M{"email": body.Email}).Decode(&userModel)
	if err == nil {
		ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already in use",
		})
		return nil
	}

	// Check if passwords match
	if body.Password != body.ConfirmPassword {
		ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Passwords do not match",
		})
		return nil
	}

	// Hash password
	argon2id := utils.NewArgon2ID()
	hashedPassword, err := argon2id.Hash(body.Password)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}

	// Create user
	userModel = &models.User{
		Username:  body.Username,
		Email:     body.Email,
		Password:  hashedPassword,
		Role:      enum.User,
		CreatedAt: time.Now().Unix(),
	}

	result, err := userCollection.InsertOne(ctx.Context(), userModel)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}
	userModel.Id = result.InsertedID.(primitive.ObjectID).Hex()

	// Send confirmation email
	go utils.ConfirmEmail(userModel)

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created",
	})
	return nil
}
