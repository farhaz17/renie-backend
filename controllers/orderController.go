package controllers

import (
	"fmt"
	"renie-backend/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetOrder(ctx *fiber.Ctx) error {
	// Parse the order ID from the URL parameters
	orderID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	orderID32 := int32(orderID)

	// Call the service layer to fetch the order
	order, err := services.GetOrderByID(orderID32)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch order")
	}

	// Return the order data as a JSON response
	return ctx.JSON(order)
}

func CreateOrder(ctx *fiber.Ctx) error {
	var req struct {
		OrderType  string `json:"order_type"`
		CustomerID int    `json:"customer_id"`
		ProductID  int    `json:"product_id"` // UUID as string
		Quantity   int    `json:"quantity"`
	}

	// Parse the incoming JSON body
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Call the service layer to create the order
	order, err := services.CreateOrder(req.OrderType, req.CustomerID, req.ProductID, req.Quantity)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
	}

	// Return the created order data as a JSON response
	return ctx.Status(fiber.StatusCreated).JSON(order)
}

// ApproveOrder handles the approval of an order
func ApproveOrder(ctx *fiber.Ctx) error {
	// Get order ID from URL params
	orderID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	// Convert to int32
	orderID32 := int32(orderID)

	// Approve the order using the service
	approvedOrder, err := services.ApproveOrder(orderID32)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to approve order")
	}

	// Return the approved order
	return ctx.JSON(approvedOrder)
}

// DispatchOrder handles the dispatch request
func DispatchOrder(ctx *fiber.Ctx) error {
	// Parse order ID
	orderIDStr := ctx.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	// Call service to dispatch order
	err = services.DispatchOrder(int32(orderID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "Order dispatched successfully"})
}

// Mark order as "Out for Delivery"
func MarkOrderOutForDelivery(ctx *fiber.Ctx) error {
	orderIDStr := ctx.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	err = services.MarkOrderOutForDelivery(int32(orderID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "Order marked as 'Out for Delivery'"})
}

// Mark order as "Delivered"
func MarkOrderDelivered(ctx *fiber.Ctx) error {
	orderIDStr := ctx.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	err = services.MarkOrderDelivered(int32(orderID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "Order marked as 'Delivered'"})
}

// Mark order as "Returned"
func MarkOrderReturned(ctx *fiber.Ctx) error {
	orderIDStr := ctx.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid order ID")
	}

	err = services.MarkOrderReturned(int32(orderID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "Order marked as 'Returned'"})
}
