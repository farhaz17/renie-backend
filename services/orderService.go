package services

import (
	"context"
	"database/sql"
	"fmt"
	"renie-backend/config"
	"renie-backend/db/sqlc"
)

func GetOrderByID(orderID int32) (*sqlc.Order, error) {
	// Use the globally initialized database connection to fetch the order
	q := sqlc.New(config.DB)
	order, err := q.GetOrderById(context.Background(), orderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching order: %w", err)
	}

	return &order, nil
}

func CreateOrder(orderType string, customerID int, productID int, quantity int) (*sqlc.CreateOrderRow, error) {
	// Initialize SQLC client
	q := sqlc.New(config.DB)

	// Call the SQLC-generated function to insert the order
	newOrder, err := q.CreateOrder(context.Background(), sqlc.CreateOrderParams{
		CustomerID: int32(customerID),
		OrderType:  sql.NullString{String: orderType, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	_, err = q.CreateOrderItem(context.Background(), sqlc.CreateOrderItemParams{
		OrderID:   newOrder.ID, // Order ID from the created order
		ProductID: int32(productID),
		Quantity:  int32(quantity),
	})

	if err != nil {
		return nil, fmt.Errorf("error creating order item: %w", err)
	}

	return &newOrder, nil
}

func ApproveOrder(orderID int32) (*sqlc.Order, error) {
	q := sqlc.New(config.DB)

	// Approve the order
	approvedOrder, err := q.ApproveOrder(context.Background(), orderID)
	if err != nil {
		return nil, fmt.Errorf("error approving order: %w", err)
	}

	return &approvedOrder, nil
}

// DispatchOrder handles order dispatching with stock validation using a transaction
func DispatchOrder(orderID int32) error {
	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	q := sqlc.New(tx) // Use transaction for queries

	// Fetch order items
	orderItems, err := q.GetOrderItemsByOrderID(context.Background(), orderID)
	if err != nil {
		return fmt.Errorf("error fetching order items: %w", err)
	}

	// Validate stock availability for each product
	for _, item := range orderItems {
		productStock, err := q.GetProductStock(context.Background(), item.ProductID)
		if err != nil {
			return fmt.Errorf("error fetching product stock: %w", err)
		}

		if productStock < item.Quantity {
			return fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}
	}

	// Deduct stock from products
	for _, item := range orderItems {
		err := q.UpdateProductStock(context.Background(), sqlc.UpdateProductStockParams{
			Stock: item.Quantity,
			ID:    item.ProductID,
		})
		if err != nil {
			return fmt.Errorf("error updating product stock: %w", err)
		}
	}

	// Update order status to "Dispatched"
	err = q.DispatchOrder(context.Background(), orderID)
	if err != nil {
		return fmt.Errorf("error updating order status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// MarkOrderOutForDelivery updates the order status to "Out for Delivery"
func MarkOrderOutForDelivery(orderID int32) error {
	q := sqlc.New(config.DB)

	err := q.MarkOrderOutForDelivery(context.Background(), orderID)
	if err != nil {
		return fmt.Errorf("error updating order status: %w", err)
	}
	return nil
}

// MarkOrderDelivered updates the order status to "Delivered"
func MarkOrderDelivered(orderID int32) error {
	q := sqlc.New(config.DB)

	err := q.MarkOrderDelivered(context.Background(), orderID)
	if err != nil {
		return fmt.Errorf("error updating order status: %w", err)
	}
	return nil
}

// MarkOrderReturned updates the order status to "Returned"
func MarkOrderReturned(orderID int32) error {
	q := sqlc.New(config.DB)

	err := q.MarkOrderReturned(context.Background(), orderID)
	if err != nil {
		return fmt.Errorf("error updating order status: %w", err)
	}
	return nil
}
