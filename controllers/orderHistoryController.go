package controllers

import (
	"jwt/database"
	"jwt/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func OrderHistory(c *fiber.Ctx) error {
	var orders []models.Order

	// Fetch all orders for the logged-in user
	tokenStr := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	userId, _ := strconv.Atoi(claims.Issuer)

	// Fetch user's orders
	if err := database.DB.Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch orders"})
	}

	return c.JSON(orders)
}
