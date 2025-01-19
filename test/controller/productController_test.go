package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"renie-backend/middlewares"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoleMiddleware_StaffCannotCreateProduct(t *testing.T) {
	// Setup Fiber app
	app := fiber.New()

	// JWT middleware (Fiber's JWT parser)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("test_secret"), // ✅ Corrected signing key
		ContextKey: "user",                // Store parsed JWT in Locals["user"]
	}))

	// Define a protected route (only admin/manager allowed)
	app.Post("/api/products", middlewares.RoleRequired("admin", "manager"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Product created"})
	})

	// Generate JWT token for staff
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  123,
		"username": "staff_user",
		"role":     "staff", // Staff role
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	signedToken, _ := token.SignedString([]byte("test_secret"))

	// Create test request
	req := httptest.NewRequest(http.MethodPost, "/api/products", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)

	// Perform test request
	resp, _ := app.Test(req)

	// Assertions
	assert.Equal(t, http.StatusForbidden, resp.StatusCode) // Expect 403 Forbidden
}
