package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Mock service function for creating orders
var mockCreateOrderService = func(orderType string, customerID, productID, quantity int) (map[string]interface{}, error) {
	if orderType == "" || customerID <= 0 || productID <= 0 || quantity <= 0 {
		return nil, errors.New("Invalid order data")
	}
	return map[string]interface{}{
		"order_id":    1,
		"order_type":  orderType,
		"customer_id": customerID,
		"product_id":  productID,
		"quantity":    quantity,
	}, nil
}

// Test handler for CreateOrder
func TestCreateOrder(t *testing.T) {
	app := fiber.New()
	app.Post("/orders", func(c *fiber.Ctx) error {
		var req struct {
			OrderType  string `json:"order_type"`
			CustomerID int    `json:"customer_id"`
			ProductID  int    `json:"product_id"`
			Quantity   int    `json:"quantity"`
		}

		// Parse the JSON request
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Call the mock service
		order, err := mockCreateOrderService(req.OrderType, req.CustomerID, req.ProductID, req.Quantity)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
		}

		return c.Status(fiber.StatusCreated).JSON(order)
	})

	// Valid request
	validOrder := map[string]interface{}{
		"order_type":  "normal",
		"customer_id": 1,
		"product_id":  10,
		"quantity":    2,
	}
	body, _ := json.Marshal(validOrder)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Invalid request (missing fields)
	invalidOrder := map[string]interface{}{
		"order_type":  "",
		"customer_id": 0,
		"product_id":  -1,
		"quantity":    0,
	}
	bodyInvalid, _ := json.Marshal(invalidOrder)

	reqInvalid := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(bodyInvalid))
	reqInvalid.Header.Set("Content-Type", "application/json")
	respInvalid, _ := app.Test(reqInvalid)
	assert.Equal(t, fiber.StatusInternalServerError, respInvalid.StatusCode)
}

// Mock service function
var mockService = func(orderID int32) error {
	if orderID <= 0 {
		return errors.New("Invalid order ID")
	}
	return nil
}

// Test handler
func TestMarkOrderOutForDelivery(t *testing.T) {
	app := fiber.New()
	app.Post("/orders/:id/out-for-delivery", func(c *fiber.Ctx) error {
		orderIDStr := c.Params("id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
		}

		err = mockService(int32(orderID))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{"message": "Order marked as 'Out for Delivery'"})
	})

	// Valid request
	req := httptest.NewRequest(http.MethodPost, "/orders/1/out-for-delivery", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Invalid request (non-numeric ID)
	reqInvalid := httptest.NewRequest(http.MethodPost, "/orders/abc/out-for-delivery", nil)
	respInvalid, _ := app.Test(reqInvalid)
	assert.Equal(t, fiber.StatusBadRequest, respInvalid.StatusCode)
}
