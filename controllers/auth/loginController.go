package auth

import (
	"time"

	"github.com/EloToJaa/quizer/db"
	"github.com/EloToJaa/quizer/enum"
	"github.com/EloToJaa/quizer/jwt"
	"github.com/EloToJaa/quizer/models"
	"github.com/EloToJaa/quizer/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

type LoginForm struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required,min=3,max=64"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
}

func LoginController(ctx *fiber.Ctx) error {
	// Get body from request; make username and email optional
	var body LoginForm

	// Parse body into struct
	if err := ctx.BodyParser(&body); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
		})
		return nil
	}

	// Validate body
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  utils.FormatValidationErrors(err),
		})
	}

	userCollection := db.GetCollection(enum.Users)
	userModel := &models.User{}

	// Check if user exists
	err := userCollection.FindOne(ctx.Context(), bson.M{"$or": []bson.M{
		{"username": body.UsernameOrEmail},
		{"email": body.UsernameOrEmail},
	}}).Decode(&userModel)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Wrong username or password",
		})
	}

	// Check if password is correct
	argon2id := utils.NewArgon2ID()
	if ok, err := argon2id.Verify(body.Password, userModel.Password); !ok || err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Wrong username or password",
		})
	}

	// Check if user is verified
	if !userModel.Verified {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User not verified",
		})
	}

	// Generate JWT
	user := utils.MapStructs(userModel, &jwt.User{}).(*jwt.User)
	accessTokenData := &jwt.TokenData{
		User:      user,
		ExpiresAt: jwt.GetAccessTokenExpirationTime().Unix(),
	}
	accessToken, err := accessTokenData.GenerateToken(jwt.GetAccessTokenSecret())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}

	refreshTokenData := &jwt.TokenData{
		User:      user,
		ExpiresAt: jwt.GetRefreshTokenExpirationTime().Unix(),
	}
	refreshToken, err := refreshTokenData.GenerateToken(jwt.GetRefreshTokenSecret(userModel.Password))
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}

	// Add refresh token to database
	refreshTokenCollection := db.GetCollection(enum.RefreshTokens)
	_, err = refreshTokenCollection.InsertOne(ctx.Context(), &models.RefreshToken{
		UserId:       userModel.Id,
		Token:        refreshToken,
		UserPassword: userModel.Password,
		CreatedAt:    time.Now().Unix(),
		ExpiresAt:    jwt.GetRefreshTokenExpirationTime().Unix(),
	})
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
		return nil
	}

	// Set cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  jwt.GetRefreshTokenExpirationTime(),
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

	// Send response
	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Success",
		"accessToken": accessToken,
	})
	return nil
}
