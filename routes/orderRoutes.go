package routes

import (
	"jwt/controllers"

	"github.com/gofiber/fiber/v2"
)

func order(app fiber.Router) {
	app.Get("/orders", controllers.GetOrders)           // Fetch all orders for the logged-in user
	app.Post("/orders", controllers.PlaceOrder)         // Place a new order
	app.Post("/orders/:id/stop", controllers.StopOrder) // Stop an open order
}
