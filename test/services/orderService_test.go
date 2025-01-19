package services_test

import (
	"context"
	"database/sql"
	"renie-backend/config"
	"renie-backend/services"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"renie-backend/db/sqlc"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func TestGetOrderByID(t *testing.T) {
	// Setup the mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize the queries with the mock database
	q := sqlc.New(db)

	// Define test data
	orderID := int32(1)
	expectedOrder := sqlc.Order{
		ID:         orderID,
		CustomerID: 2,
		OrderType:  sql.NullString{String: "Normal", Valid: true},
		Status:     sql.NullString{String: "Created", Valid: true},
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true}, // ✅ Correct Fix
		UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
	}

	// Mock SQL query
	mock.ExpectQuery(`SELECT id, order_type, customer_id, status, created_at, updated_at FROM orders WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_type", "customer_id", "status", "created_at", "updated_at"}).
			AddRow(
				expectedOrder.ID,
				expectedOrder.OrderType.String, // Use .String for sql.NullString
				expectedOrder.CustomerID,
				expectedOrder.Status.String,  // Use .String for sql.NullString
				expectedOrder.CreatedAt.Time, // Use .Time for sql.NullTime
				expectedOrder.UpdatedAt.Time, // Use .Time for sql.NullTime
			))

	// Call the function
	order, err := q.GetOrderById(context.Background(), orderID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID, order.ID)
	assert.Equal(t, expectedOrder.CustomerID, order.CustomerID)
	assert.Equal(t, expectedOrder.OrderType, order.OrderType)
	assert.Equal(t, expectedOrder.Status, order.Status)

	// Verify all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateOrder(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize SQLC Queries with mock DB
	q := sqlc.New(db)

	// Define test data
	customerID := int32(2)
	orderType := "Normal"
	productID := int32(3)
	quantity := int32(5)
	orderID := int32(1)
	createdAt := time.Now()

	// Mock the CreateOrder query (Fixed SQL regex)
	mock.ExpectQuery(`INSERT INTO orders \(customer_id, order_type, status\) VALUES \(\$1, \$2, 'Created'\) RETURNING id, customer_id, order_type, status, created_at`).
		WithArgs(customerID, orderType).
		WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "order_type", "status", "created_at"}).
			AddRow(orderID, customerID, orderType, "Created", createdAt))

	// Mock the CreateOrderItem query
	mock.ExpectQuery(`INSERT INTO order_items \(order_id, product_id, quantity\) VALUES \(\$1, \$2, \$3\) RETURNING id, order_id, product_id, quantity, created_at, updated_at`).
		WithArgs(orderID, productID, quantity).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "created_at", "updated_at"}).
			AddRow(1, orderID, productID, quantity, createdAt, createdAt))

	// Call the function to test
	newOrder, err := q.CreateOrder(context.Background(), sqlc.CreateOrderParams{
		CustomerID: customerID,
		OrderType:  sql.NullString{String: orderType, Valid: true},
	})
	require.NoError(t, err) // Ensure no error occurs

	_, err = q.CreateOrderItem(context.Background(), sqlc.CreateOrderItemParams{
		OrderID:   newOrder.ID,
		ProductID: productID,
		Quantity:  quantity,
	})
	require.NoError(t, err) // Ensure no error occurs

	// Assertions
	assert.Equal(t, orderID, newOrder.ID)
	assert.Equal(t, customerID, newOrder.CustomerID)
	assert.Equal(t, orderType, newOrder.OrderType.String)

	// Verify all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDispatchOrder(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Temporarily replace the global DB connection
	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }() // Restore original DB after test

	// Test data
	orderID := int32(1)
	productID := int32(2)
	quantity := int32(5)
	stock := int32(10)

	// Mock GetOrderItemsByOrderID query
	mock.ExpectQuery(`SELECT order_id, product_id, quantity FROM order_items WHERE order_id = \$1`).
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "product_id", "quantity"}).
			AddRow(orderID, productID, quantity))

	// Mock GetProductStock query
	mock.ExpectQuery(`SELECT stock FROM products WHERE id = \$1`).
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"stock"}).
			AddRow(stock))

	// Mock UpdateProductStock query
	mock.ExpectExec(`UPDATE products SET stock = stock - \$1 WHERE id = \$2`).
		WithArgs(quantity, productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock DispatchOrder query
	mock.ExpectExec(`UPDATE orders SET status = 'Dispatched' WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function to test (without modifying DispatchOrder)
	err = services.DispatchOrder(orderID)

	// Assertions
	assert.NoError(t, err)

	// Verify all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test MarkOrderOutForDelivery
func TestMarkOrderOutForDelivery(t *testing.T) {
	db, mock := setupMockDB()
	defer db.Close()

	q := sqlc.New(db)

	orderID := int32(1)
	mock.ExpectExec("UPDATE orders SET status = 'Out for Delivery' WHERE id = \\$1").
		WithArgs(orderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := q.MarkOrderOutForDelivery(context.Background(), orderID)
	assert.NoError(t, err)
}
