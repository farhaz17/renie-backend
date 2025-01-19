package controllers

import (
	"fmt"
	"renie-backend/services"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ProductRequest represents the JSON body for creating a product
type ProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
	Stock       int32  `json:"stock"`
}

// CreateProduct handles the API request for adding a new product
func CreateProduct(ctx *fiber.Ctx) error {
	var req ProductRequest

	// Parse the incoming JSON request
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Call the service layer to create the product
	product, err := services.CreateProduct(req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create product: %v", err))
	}

	// Return the created product as a JSON response
	return ctx.Status(fiber.StatusCreated).JSON(product)
}

func GetProduct(ctx *fiber.Ctx) error {
	// Parse the product ID from the URL parameters
	productIDStr := ctx.Params("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid product ID")
	}

	// Call the service layer to fetch the product by ID
	product, err := services.GetProductByID(int32(productID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch product: %v", err))
	}

	// Return the product data as a JSON response
	return ctx.JSON(product)
}
