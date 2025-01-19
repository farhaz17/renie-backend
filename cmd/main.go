package main

import (
	"log"
	"renie-backend/config"
	"renie-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.Init()

	routes.SetupRoutes(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
