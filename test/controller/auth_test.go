package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Setup Fiber app with the /orders/:id route
func setupApp() *fiber.App {
	app := fiber.New()

	// Mock protected route that requires authentication
	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.JSON(fiber.Map{"message": "Order details"})
	})

	return app
}

// Test accessing /orders/:id without authentication (should return 401)
func TestUnauthorizedOrderAccess(t *testing.T) {
	app := setupApp()

	// Make a request to GET /orders/1 without Authorization header
	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	resp, _ := app.Test(req)

	// Assert that the response status is 401 Unauthorized
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	// Assert that the response contains the error message
	expectedResponse := `{"error":"Unauthorized"}`
	buf := make([]byte, len(expectedResponse))
	resp.Body.Read(buf) // Read response body
	assert.JSONEq(t, expectedResponse, string(buf))
}
