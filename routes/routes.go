package routes

import (
	"renie-backend/controllers"

	"renie-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Login route (no protection)
	app.Post("/login", controllers.Login)

	// Protected API group
	api := app.Group("/api", middlewares.JWTProtected())

	// Order routes
	api.Get("/orders/:id", middlewares.RoleRequired("admin", "manager", "staff"), controllers.GetOrder)
	api.Post("/orders", middlewares.RoleRequired("admin", "manager", "staff"), controllers.CreateOrder)
	api.Put("/orders/:id/approve", middlewares.RoleRequired("admin", "manager"), controllers.ApproveOrder)
	api.Put("/orders/:id/dispatch", middlewares.RoleRequired("admin", "manager"), controllers.DispatchOrder)
	api.Put("/orders/:id/out-for-delivery", middlewares.RoleRequired("admin", "manager"), controllers.MarkOrderOutForDelivery)
	api.Put("/orders/:id/delivered", middlewares.RoleRequired("admin", "manager"), controllers.MarkOrderDelivered)
	api.Put("/orders/:id/returned", middlewares.RoleRequired("admin", "manager"), controllers.MarkOrderReturned)

	// Product routes
	api.Post("/products", middlewares.RoleRequired("admin", "manager"), controllers.CreateProduct)
	api.Get("/products/:id", middlewares.RoleRequired("admin", "manager"), controllers.GetProduct)
}
