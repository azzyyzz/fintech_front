package controllers

import (
	"jwt/database"
	"jwt/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Deposit(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	amount, err := strconv.ParseFloat(c.Query("amount"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid amount"})
	}

	// Assume user is authenticated and the user ID is stored in the JWT claim
	tokenStr := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	userId, _ := strconv.Atoi(claims.Issuer)

	// Deposit the amount into the user's wallet
	err = models.Deposit(database.DB, uint(userId), symbol, amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to deposit"})
	}

	return c.JSON(fiber.Map{"message": "Deposit successful"})
}
