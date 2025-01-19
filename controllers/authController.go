package controllers

import (
	"context"
	"os"
	"time"

	"renie-backend/config"
	"renie-backend/db/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the expected login body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login generates a JWT token for authenticated users
func Login(ctx *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	q := sqlc.New(config.DB)

	// Fetch user from DB
	user, err := q.GetUserByUsername(context.Background(), req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	// Return JWT token
	return ctx.JSON(fiber.Map{"token": signedToken})
}
