package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import PostgreSQL driver
	"github.com/joho/godotenv"
)

// Global database instance
var DB *sql.DB

// Init initializes environment variables and database connection
func Init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Construct DB connection string
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Connect to database
	DB, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("✅ Database connected successfully")
}
