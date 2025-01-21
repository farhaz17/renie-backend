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

