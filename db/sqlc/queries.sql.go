// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package sqlc

import (
	"context"
	"database/sql"
)

const approveOrder = `-- name: ApproveOrder :one
UPDATE orders
SET status = 'Approved', updated_at = NOW()
WHERE id = $1
RETURNING id, order_type, customer_id, status, created_at, updated_at
`

func (q *Queries) ApproveOrder(ctx context.Context, id int32) (Order, error) {
	row := q.db.QueryRowContext(ctx, approveOrder, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.OrderType,
		&i.CustomerID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (customer_id, order_type, status) VALUES ($1, $2, 'Created') 
RETURNING id, customer_id, order_type, status, created_at
`

type CreateOrderParams struct {
	CustomerID int32          `json:"customer_id"`
	OrderType  sql.NullString `json:"order_type"`
}

type CreateOrderRow struct {
	ID         int32          `json:"id"`
	CustomerID int32          `json:"customer_id"`
	OrderType  sql.NullString `json:"order_type"`
	Status     sql.NullString `json:"status"`
	CreatedAt  sql.NullTime   `json:"created_at"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (CreateOrderRow, error) {
	row := q.db.QueryRowContext(ctx, createOrder, arg.CustomerID, arg.OrderType)
	var i CreateOrderRow
	err := row.Scan(
		&i.ID,
		&i.CustomerID,
		&i.OrderType,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const createOrderItem = `-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity)
VALUES ($1, $2, $3)
RETURNING id, order_id, product_id, quantity, created_at, updated_at
`

type CreateOrderItemParams struct {
	OrderID   int32 `json:"order_id"`
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, createOrderItem, arg.OrderID, arg.ProductID, arg.Quantity)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, description, price, stock, created_at, updated_at
`

type CreateProductParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Price       int32          `json:"price"`
	Stock       int32          `json:"stock"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Stock,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Stock,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const dispatchOrder = `-- name: DispatchOrder :exec
UPDATE orders SET status = 'Dispatched' WHERE id = $1
`

func (q *Queries) DispatchOrder(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, dispatchOrder, id)
	return err
}

const getOrderById = `-- name: GetOrderById :one
SELECT id, order_type, customer_id, status, created_at, updated_at FROM orders WHERE id = $1
`

func (q *Queries) GetOrderById(ctx context.Context, id int32) (Order, error) {
	row := q.db.QueryRowContext(ctx, getOrderById, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.OrderType,
		&i.CustomerID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrderItemsByOrderID = `-- name: GetOrderItemsByOrderID :many
SELECT order_id, product_id, quantity
FROM order_items
WHERE order_id = $1
`

type GetOrderItemsByOrderIDRow struct {
	OrderID   int32 `json:"order_id"`
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (q *Queries) GetOrderItemsByOrderID(ctx context.Context, orderID int32) ([]GetOrderItemsByOrderIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrderItemsByOrderID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderItemsByOrderIDRow
	for rows.Next() {
		var i GetOrderItemsByOrderIDRow
		if err := rows.Scan(&i.OrderID, &i.ProductID, &i.Quantity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductByID = `-- name: GetProductByID :one
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE id = $1
`

func (q *Queries) GetProductByID(ctx context.Context, id int32) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProductByID, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Stock,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductStock = `-- name: GetProductStock :one
SELECT stock FROM products WHERE id = $1
`

func (q *Queries) GetProductStock(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRowContext(ctx, getProductStock, id)
	var stock int32
	err := row.Scan(&stock)
	return stock, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, password, role
FROM users
WHERE username = $1
`

type GetUserByUsernameRow struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i GetUserByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const markOrderDelivered = `-- name: MarkOrderDelivered :exec
UPDATE orders SET status = 'Delivered' WHERE id = $1
`

func (q *Queries) MarkOrderDelivered(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, markOrderDelivered, id)
	return err
}

const markOrderOutForDelivery = `-- name: MarkOrderOutForDelivery :exec
UPDATE orders SET status = 'Out for Delivery' WHERE id = $1
`

func (q *Queries) MarkOrderOutForDelivery(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, markOrderOutForDelivery, id)
	return err
}

const markOrderReturned = `-- name: MarkOrderReturned :exec
UPDATE orders SET status = 'Returned' WHERE id = $1
`

func (q *Queries) MarkOrderReturned(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, markOrderReturned, id)
	return err
}

const updateOrderStatus = `-- name: UpdateOrderStatus :exec
UPDATE orders SET status = $1 WHERE id = $2
`

type UpdateOrderStatusParams struct {
	Status sql.NullString `json:"status"`
	ID     int32          `json:"id"`
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateOrderStatus, arg.Status, arg.ID)
	return err
}

const updateProductStock = `-- name: UpdateProductStock :exec
UPDATE products SET stock = stock - $1 WHERE id = $2
`

type UpdateProductStockParams struct {
	Stock int32 `json:"stock"`
	ID    int32 `json:"id"`
}

func (q *Queries) UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) error {
	_, err := q.db.ExecContext(ctx, updateProductStock, arg.Stock, arg.ID)
	return err
}
