# Running the Project Locally

## Prerequisites

Ensure you have the following installed:

- **Go 1.20+**: [Download and install Go](https://golang.org/dl/)
- **PostgreSQL**: Install and set up PostgreSQL
- **SQLC**: Install using the command:
  ```sh
  go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
  ```

## Clone the Repository

```sh
git clone https://github.com/your-repo/renie-backend.git
cd renie-backend
```

## Install Dependencies

```sh
go mod tidy
```

## Setup Environment Variables

Create a `.env` file in the root directory and configure the database connection:

```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=renie_db

JWT_SECRET=your_secret_key - (generated online)
```

## Run Migrations

Apply database migrations using the following command:

```sh
migrate -path ./db/migrations -database "postgres://postgres:postgres@localhost:5432/renie_db?sslmode=disable" -verbose up
```

## Run Unit Tests

```sh
go test ./test/services ./test/controller -v
```

## Start the Server

Run the server with:

```sh
go run cmd/main.go
```

Your server should now be running and ready for use!


# API Documentation

## Authentication

### Login
**Endpoint:** `POST /login`
**Description:** Authenticate a user and generate a JWT token.

#### Request Body (JSON):
```json
{
    "username": "admin",
    "password": "password" 
}
```

#### Response (JSON):
```json
{
  "token": "your_jwt_token_here"
}
```

## Orders

### Get Order by ID
**Endpoint:** `GET /api/orders/:id`
**Roles Required:** `admin`, `manager`, `staff`
**Description:** Retrieve details of a specific order.

#### Response (JSON):
```json
{
    "id": 1,
    "order_type": {
        "String": "Normal",
        "Valid": true
    },
    "customer_id": 1,
    "status": {
        "String": "Created",
        "Valid": true
    },
    "created_at": {
        "Time": "2025-01-19T08:23:04.93885Z",
        "Valid": true
    },
    "updated_at": {
        "Time": "2025-01-19T08:23:04.93885Z",
        "Valid": true
    }
}
```

### Create Order
**Endpoint:** `POST /api/orders`
**Roles Required:** `admin`, `manager`, `staff`
**Description:** Create a new order.

#### Request Body (JSON):
```json
{
    "order_type" : "Normal",
    "customer_id" : 1,
    "product_id": 1,
    "quantity": 5
}
```

#### Response (JSON):
```json
{
    "id": 2,
    "customer_id": 1,
    "order_type": {
        "String": "Normal",
        "Valid": true
    },
    "status": {
        "String": "Created",
        "Valid": true
    },
    "created_at": {
        "Time": "2025-01-19T23:14:40.858951Z",
        "Valid": true
    }
}
```

### Approve Order
**Endpoint:** `PUT /api/orders/:id/approve`
**Roles Required:** `admin`, `manager`
**Description:** Approve an order for processing.

#### Response (JSON):
```json
{
  "message": "Order approved successfully"
}
```

### Dispatch Order
**Endpoint:** `PUT /api/orders/:id/dispatch`
**Roles Required:** `admin`, `manager`
**Description:** Dispatch an approved order.

#### Response (JSON):
```json
{
  "message": "Order dispatched successfully"
}
```

### Mark Order Out for Delivery
**Endpoint:** `PUT /api/orders/:id/out-for-delivery`
**Roles Required:** `admin`, `manager`
**Description:** Mark an order as "Out for Delivery".

#### Response (JSON):
```json
{
  "message": "Order marked as 'Out for Delivery'"
}
```

### Mark Order as Delivered
**Endpoint:** `PUT /api/orders/:id/delivered`
**Roles Required:** `admin`, `manager`
**Description:** Update order status to "Delivered".

#### Response (JSON):
```json
{
  "message": "Order delivered successfully"
}
```

### Mark Order as Returned
**Endpoint:** `PUT /api/orders/:id/returned`
**Roles Required:** `admin`, `manager`
**Description:** Mark an order as "Returned".

#### Response (JSON):
```json
{
  "message": "Order marked as returned"
}
```

## Products

### Create Product
**Endpoint:** `POST /api/products`
**Roles Required:** `admin`, `manager`
**Description:** Create a new product.

#### Request Body (JSON):
```json
{
    "name": "Smart Bin X3",
    "description": "An advanced smart waste bin with AI sorting.",
    "price": 199,
    "stock": 50
}
```

#### Response (JSON):
```json
{
    "id": 4,
    "name": "Smart Bin X3",
    "description": {
        "String": "An advanced smart waste bin with AI sorting.",
        "Valid": true
    },
    "price": 199,
    "stock": 50,
    "created_at": {
        "Time": "2025-01-19T08:39:47.663534Z",
        "Valid": true
    },
    "updated_at": {
        "Time": "2025-01-19T08:39:47.663534Z",
        "Valid": true
    }
}
```

### Get Product by ID
**Endpoint:** `GET /api/products/:id`
**Roles Required:** `admin`, `manager`
**Description:** Retrieve product details.

#### Response (JSON):
```json
{
    "id": 1,
    "name": "Smart Bin X1",
    "description": {
        "String": "An advanced smart waste bin with AI sorting.",
        "Valid": true
    },
    "price": 199,
    "stock": 50,
    "created_at": {
        "Time": "2025-01-19T05:46:32.079408Z",
        "Valid": true
    },
    "updated_at": {
        "Time": "2025-01-19T05:46:32.079408Z",
        "Valid": true
    }
}
```

# Unit Test Documentation

This document provides an overview of the unit tests implemented for the application, focusing on critical functionalities such as order management and role-based access control.

## Test Cases

### TestUnauthorizedOrderAccess

**Test Logic:**
- Simulate a request to create or view an order without valid authentication.
- Assert that the response is a `401 Unauthorized` error.

### TestCreateOrder

**Test Logic:**
- Assert that the order is successfully created with a `201 Created` response.

### TestMarkOrderOutForDelivery

**Test Logic:**
- Assert that the response has a `200 OK` status.

### TestRoleMiddleware_StaffCannotCreateProduct

**Test Logic:**
- Simulate a `POST /products` request with a "staff" role user.
- Assert that the response is `403 Forbidden`.

### TestGetOrderByID

**Test Logic:**
- Mock the database query for retrieving an order.
- Assert that the correct order is returned.

### TestCreateOrder (Service)

**Test Logic:**
- Simulate valid order creation input.
- Assert that a new order is added to the database and the correct data is returned.

### TestDispatchOrder

**Test Logic:**
- Simulate an existing order in the database.
- Call the service method to dispatch the order.
- Assert that the order’s status is updated to "dispatched".

### TestMarkOrderOutForDelivery (Service)

**Test Logic:**
- Simulate an existing order in the database.
- Assert that the order’s status is updated accordingly.

## Test Execution

Each of these tests ensures the application behaves correctly under various scenarios.
To run the tests, use the following command:

```bash
go test ./test/services ./test/controller -v




