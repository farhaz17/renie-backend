-- name: CreateOrder :one
INSERT INTO orders (customer_id, order_type, status) VALUES ($1, $2, 'Created') 
RETURNING id, customer_id, order_type, status, created_at;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity)
VALUES ($1, $2, $3)
RETURNING id, order_id, product_id, quantity, created_at, updated_at;

-- name: UpdateOrderStatus :exec
UPDATE orders SET status = $1 WHERE id = $2;

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = $1;

-- name: ApproveOrder :one
UPDATE orders
SET status = 'Approved', updated_at = NOW()
WHERE id = $1
RETURNING id, order_type, customer_id, status, created_at, updated_at;

-- name: GetOrderItemsByOrderID :many
SELECT order_id, product_id, quantity
FROM order_items
WHERE order_id = $1;

-- name: GetProductStock :one
SELECT stock FROM products WHERE id = $1;

-- name: UpdateProductStock :exec
UPDATE products SET stock = stock - $1 WHERE id = $2;

-- name: DispatchOrder :exec
UPDATE orders SET status = 'Dispatched' WHERE id = $1;


-- name: MarkOrderOutForDelivery :exec
UPDATE orders SET status = 'Out for Delivery' WHERE id = $1;

-- name: MarkOrderDelivered :exec
UPDATE orders SET status = 'Delivered' WHERE id = $1;

-- name: MarkOrderReturned :exec
UPDATE orders SET status = 'Returned' WHERE id = $1;






-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, description, price, stock, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE id = $1;



-- name: GetUserByUsername :one
SELECT id, username, email, password, role
FROM users
WHERE username = $1;
