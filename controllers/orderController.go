package controllers

import (
	"fmt"
	"jwt/database"
	"jwt/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetOrders(c *fiber.Ctx) error {
	// Get user from JWT cookie
	tokenStr := c.Cookies("jwt")
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

	// Fetch orders for the user, sorted by status
	var orders []models.Order
	fmt.Println(userId)
	database.DB.Where("user_id = ?", userId).Order("type ASC").Find(&orders)

	return c.JSON(orders)
}

type OrderRequest struct {
	Symbol  string  `json:"symbol"`
	SellBuy string  `json:"sellbuy"` // "buy" or "sell"
	Price   float64 `json:"price"`
	Amount  float64 `json:"amount"`
}

func PlaceOrder(c *fiber.Ctx) error {
	var orderRequest struct {
		UserId    uint    `json:"user_id"`
		Price     float64 `json:"price"`
		Amount    float64 `json:"amount"`
		Symbol    string  `json:"symbol"`
		OrderType string  `json:"type"`
		SellBuy   string  `json:"sellbuy"`
	}

	fmt.Println(c.BodyParser((&orderRequest)))
	// Parse the JSON request body
	if err := c.BodyParser(&orderRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Create the order using the CreateOrder function
	order, err := models.CreateOrder(database.DB, orderRequest.UserId, orderRequest.Price, orderRequest.Amount, orderRequest.Symbol, orderRequest.OrderType, orderRequest.SellBuy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the created order as a response
	return c.Status(fiber.StatusCreated).JSON(order)
}

func StopOrder(c *fiber.Ctx) error {
	// Get order ID from URL params
	orderID := c.Params("id")

	// Update the order status to "stopped"
	var order models.Order
	database.DB.First(&order, orderID)
	if order.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	order.Type = "stopped"
	database.DB.Save(&order)

	return c.JSON(fiber.Map{
		"message": "Order stopped successfully",
	})
}
