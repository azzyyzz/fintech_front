package controllers

import (
	"fmt"
	"jwt/database"
	"jwt/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
) // BalanceResponse struct to send balance data in response
type BalanceResponse struct {
	Usdt float64 `json:"usdt"`
	Btc  float64 `json:"btc"`
	Eth  float64 `json:"eth"`
}

func GetBalances(c *fiber.Ctx) error {
	// Get the userID from the query paramete
	fmt.Println("JUST TESTING HERE1")
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	fmt.Println("JUST TESTING HERE2")

	// Fetch the balances from the database for the user
	var wallets []models.Wallet
	if err := database.DB.Where("user_id = ?", data["userId"]).Find(&wallets).Error; err != nil {
		log.Println("Error retrieving balances:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balances",
		})
	}

	fmt.Println("JUST TESTING HERE")

	// Map the wallet balances to the response format
	balances := BalanceResponse{}
	for _, wallet := range wallets {
		switch wallet.Symbol {
		case "USDT":
			balances.Usdt = wallet.Balance
		case "BTC":
			balances.Btc = wallet.Balance
		case "ETH":
			balances.Eth = wallet.Balance
		}
	}

	// Respond with the balance data in JSON format
	return c.JSON(balances)
}

func GetWallets(c *fiber.Ctx) error {
	// Get user from JWT cookie
	tokenStr := c.Cookies("jwt")

	// Parse the token to get user ID
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	userId := claims.Issuer

	// Fetch wallets for the user
	var wallets []models.Wallet
	database.DB.Where("user_id = ?", userId).Find(&wallets)

	return c.JSON(wallets)
}

func DepositFunds(c *fiber.Ctx) error {
	// Parse request body
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	// Validate input
	amount := data["amount"].(float64)
	symbol := data["symbol"].(string)

	if amount <= 0 || (symbol != "BTC" && symbol != "ETH" && symbol != "USDT") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid amount or symbol",
		})
	}

	// Get user from JWT cookie
	tokenStr := c.Cookies("jwt")
	fmt.Print(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": symbol,
		})
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	userId := claims.Issuer

	// Update wallet balance
	var wallet models.Wallet
	database.DB.Where("user_id = ? AND symbol = ?", userId, symbol).First(&wallet)
	if wallet.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Wallet not found",
		})
	}

	wallet.Balance += amount
	database.DB.Save(&wallet)

	return c.JSON(wallet)
}
