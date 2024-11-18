package routes

import (
	"fmt"
	"jwt/database"
	"jwt/models"
	"log"
	"math"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func MatchOrders(db *gorm.DB) error {
	var buyOrders, sellOrders []models.Order

	// Fetch buy orders (Price descending order, Type not closed)
	err := db.Where("sell_buy = ? AND type != ?", "buy", "closed").
		Order("price DESC, date ASC").
		Find(&buyOrders).Error
	if err != nil {
		return fmt.Errorf("could not fetch buy orders: %v", err)
	}

	// Fetch sell orders (Price ascending order, Type not closed)
	err = db.Where("sell_buy = ? AND type != ?", "sell", "closed").
		Order("price ASC, date ASC").
		Find(&sellOrders).Error
	if err != nil {
		return fmt.Errorf("could not fetch sell orders: %v", err)
	}

	// Match orders (price-time priority)
	for _, buyOrder := range buyOrders {
		fmt.Println("TRYING2")
		fmt.Println(buyOrder)
		for _, sellOrder := range sellOrders {
			fmt.Println(sellOrder)
			// If the buy price is higher than or equal to the sell price, a match occurs
			if buyOrder.Price >= sellOrder.Price && buyOrder.UserId != sellOrder.UserId {
				// Calculate the matched amount (smallest amount between buy and sell order)
				matchAmount := math.Min(buyOrder.AmountLeft, sellOrder.AmountLeft)

				// Process the matching orders
				err := processMatch(db, buyOrder, sellOrder, matchAmount)
				if err != nil {
					return fmt.Errorf("could not process match: %v", err)
				}

				// Update the order amounts after the match
				buyOrder.AmountLeft -= matchAmount
				sellOrder.AmountLeft -= matchAmount

				// If the order is fully filled, update the order type to "closed"
				if buyOrder.AmountLeft == 0 {
					buyOrder.Type = "closed"
				}
				if sellOrder.AmountLeft == 0 {
					sellOrder.Type = "closed"
				}

				// Save the updated orders
				db.Save(&buyOrder)
				db.Save(&sellOrder)

				// Break the inner loop if the buy order is filled
				if buyOrder.AmountLeft == 0 {
					break
				}
			}
		}
	}

	fmt.Println("JUST FINISHED LOL")

	return nil
}

// processMatch updates the wallets of the users involved in the order match
func processMatch(db *gorm.DB, buyOrder models.Order, sellOrder models.Order, matchAmount float64) error {
	// Update the wallets after the match
	err := updateWallets(db, buyOrder, sellOrder, matchAmount)
	if err != nil {
		return fmt.Errorf("could not update wallets: %v", err)
	}

	// Optionally log the match
	log.Printf("Matched order: Buy Order ID %d (Price: %.2f) with Sell Order ID %d (Price: %.2f)",
		buyOrder.Id, buyOrder.Price, sellOrder.Id, sellOrder.Price)

	return nil
}

// updateWallets updates the users' wallets based on the matched orders
func updateWallets(db *gorm.DB, buyOrder models.Order, sellOrder models.Order, matchAmount float64) error {
	var buyUserWallet, sellUserWallet models.Wallet

	// Fetch the wallets for both users involved in the order matching
	err := db.Where("user_id = ? AND symbol = ?", buyOrder.UserId, getBaseCurrency(buyOrder.Symbol)).First(&buyUserWallet).Error
	if err != nil {
		return fmt.Errorf("could not fetch buy user's wallet: %v", err)
	}

	err = db.Where("user_id = ? AND symbol = ?", sellOrder.UserId, getBaseCurrency(sellOrder.Symbol)).First(&sellUserWallet).Error
	if err != nil {
		return fmt.Errorf("could not fetch sell user's wallet: %v", err)
	}

	// Update Buyer's Wallet (decrease USDT, increase BTC)
	if buyOrder.SellBuy == "buy" {
		// Buy order: decrease USDT, increase BTC
		buyUserWallet.Balance -= matchAmount * sellOrder.Price // Decrease USDT
		err = db.Save(&buyUserWallet).Error
		if err != nil {
			return fmt.Errorf("could not update buy user's wallet: %v", err)
		}

		// Increase BTC
		err = db.Where("user_id = ? AND symbol = ?", buyOrder.UserId, getBaseCurrency(buyOrder.Symbol)).First(&buyUserWallet).Error
		if err != nil {
			return fmt.Errorf("could not fetch buy user's BTC wallet: %v", err)
		}
		buyUserWallet.Balance += matchAmount // Increase BTC
		err = db.Save(&buyUserWallet).Error
		if err != nil {
			return fmt.Errorf("could not update buy user's BTC wallet: %v", err)
		}
	}

	// Update Seller's Wallet (increase USDT, decrease BTC)
	if sellOrder.SellBuy == "sell" {
		// Sell order: increase USDT, decrease BTC
		sellUserWallet.Balance += matchAmount * buyOrder.Price // Increase USDT
		err = db.Save(&sellUserWallet).Error
		if err != nil {
			return fmt.Errorf("could not update sell user's wallet: %v", err)
		}

		// Decrease BTC
		err = db.Where("user_id = ? AND symbol = ?", sellOrder.UserId, getBaseCurrency(sellOrder.Symbol)).First(&sellUserWallet).Error
		if err != nil {
			return fmt.Errorf("could not fetch sell user's BTC wallet: %v", err)
		}
		sellUserWallet.Balance -= matchAmount // Decrease BTC
		err = db.Save(&sellUserWallet).Error
		if err != nil {
			return fmt.Errorf("could not update sell user's BTC wallet: %v", err)
		}
	}

	return nil
}

// Helper function to get the base currency from the symbol (e.g., BTCUSDT -> BTC)
func getBaseCurrency(symbol string) string {
	if symbol == "BTCUSDT" || symbol == "ETHUSDT" {
		return symbol[:3] // Return "BTC" or "ETH" as base currency
	}
	return ""
}

// RunMatchEngine starts the matching engine in a separate goroutine
func RunMatchEngine(db *gorm.DB) {
	// Match orders every 2 seconds
	for {
		err := MatchOrders(db)
		if err != nil {
			log.Println("Error matching orders:", err)
		}

		// Wait for 2 seconds before the next matching cycle
		time.Sleep(2 * time.Second)
	}
}

func Run() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true, // Allow cookies to be sent
		AllowOrigins:     "*",  // Adjust this based on your frontend's origin
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// app.Use(cors.New(cors.Config{
	// 	AllowCredentials: true,
	// }))

	app.Route("/auth", auth)
	app.Route("/wallets", wallet) // Add wallet routes
	app.Route("/orders", order)   // Add order routes
	// app.Get("/justtesting", printCookies)
	go RunMatchEngine(database.DB)

	log.Fatal(app.Listen(fmt.Sprintf("localhost:%s", os.Getenv("PORT"))))

}
