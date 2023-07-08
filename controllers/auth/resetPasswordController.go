package auth

import (
	"time"

	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/models"
	"github.com/EloToJaa/quizer/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

type ResetPasswordForm struct {
	OldPassword     string `json:"oldPassword" validate:"required,min=8,max=64"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=64"`
}

func ResetPasswordController(ctx *fiber.Ctx) error {
	var body ResetPasswordForm

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

	// Check if password and confirm password match
	if body.Password != body.ConfirmPassword {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password and confirm password do not match",
		})
		return nil
	}

	id := ctx.Params("id")
	resetPasswordCollection := db.GetCollection(enum.ResetPassword)
	resetPasswordModel := &models.ResetPassword{}

	// Check if verification token exists in database
	objectId, _ := primitive.ObjectIDFromHex(id)

	err := resetPasswordCollection.FindOne(ctx.Context(), primitive.M{"_id": objectId}).Decode(resetPasswordModel)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Token not found",
		})
	}

	expired := time.Now().Unix() > resetPasswordModel.ExpiresAt
	if expired {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Token expired",
		})
	}

	// Verify user
	userCollection := db.GetCollection(enum.Users)
	userModel := &models.User{}
	objectId, _ = primitive.ObjectIDFromHex(resetPasswordModel.UserId)

	// Check if user exists
	err = userCollection.FindOne(ctx.Context(), primitive.M{"_id": objectId}).Decode(&userModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Check if old password is correct
	argon2 := utils.NewArgon2ID()
	match, err := argon2.Verify(body.OldPassword, userModel.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if !match {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Old password is incorrect",
		})
	}

	// Hash new password
	hashedPassword, err := argon2.Hash(body.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Update user password
	_, err = userCollection.UpdateOne(ctx.Context(), primitive.M{"_id": objectId}, primitive.M{"$set": primitive.M{"password": hashedPassword}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Delete all reset password tokens for user
	_, err = resetPasswordCollection.DeleteMany(ctx.Context(), primitive.M{"userId": resetPasswordModel.UserId})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated",
	})

	return nil
}
