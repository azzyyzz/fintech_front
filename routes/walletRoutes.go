package routes

import (
	"jwt/controllers"

	"github.com/gofiber/fiber/v2"
)

func wallet(app fiber.Router) {
	app.Get("/wallets", controllers.GetWallets)            // Fetch wallets for the logged-in user
	app.Post("/wallets/deposit", controllers.DepositFunds) // Deposit funds to a wallet // Fetch wallets for the logged-in user
	app.Post("/wallets/balances", controllers.GetBalances)
}
