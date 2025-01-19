package services

import (
	"context"
	"database/sql"
	"fmt"

	"renie-backend/config"
	"renie-backend/db/sqlc"
)

// CreateProduct adds a new product to the database
func CreateProduct(name, description string, price int32, stock int32) (*sqlc.Product, error) {
	q := sqlc.New(config.DB)

	// Insert the product into the database
	newProduct, err := q.CreateProduct(context.Background(), sqlc.CreateProductParams{
		Name:        name,
		Description: sql.NullString{String: description, Valid: true},
		Price:       price,
		Stock:       stock,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return &newProduct, nil
}

func GetProductByID(productID int32) (*sqlc.Product, error) {
	// Initialize the SQLC client
	q := sqlc.New(config.DB)

	// Query the database to get the product by its ID
	product, err := q.GetProductByID(context.Background(), productID)
	if err != nil {
		return nil, fmt.Errorf("error fetching product: %w", err)
	}

	return &product, nil
}
