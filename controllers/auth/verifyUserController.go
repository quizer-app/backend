package auth

import (
	"time"

	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VerifyUserController(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	verifyCollection := db.GetCollection(enum.Verify)
	verifyModel := &models.Verify{}

	// Check if verification token exists in database
	objectId, _ := primitive.ObjectIDFromHex(id)

	err := verifyCollection.FindOne(ctx.Context(), bson.M{"_id": objectId}).Decode(&verifyModel)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Token not found",
		})
	}

	expired := time.Now().Unix() > verifyModel.ExpiresAt
	if expired {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Token expired",
		})
	}

	// Verify user
	userCollection := db.GetCollection(enum.Users)
	userModel := &models.User{}
	objectId, _ = primitive.ObjectIDFromHex(verifyModel.UserId)

	// Check if user exists
	err = userCollection.FindOne(ctx.Context(), bson.M{"_id": objectId}).Decode(&userModel)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Check if user is already verified
	if userModel.Verified {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already verified",
		})
	}

	// Verify user
	userModel.Verified = true
	userModel.Id = ""

	_, err = userCollection.UpdateOne(ctx.Context(), bson.M{"_id": objectId}, bson.M{"$set": userModel})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Delete all verification tokens
	_, err = verifyCollection.DeleteMany(ctx.Context(), bson.M{"userId": verifyModel.UserId})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User verified",
	})

	return nil
}
