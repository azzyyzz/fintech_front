package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	Id            uint      `json:"id" gorm:"primaryKey"`
	UserId        uint      `json:"user_id"`
	Price         float64   `json:"price"`
	AmountInitial float64   `json:"amount_initial"`
	AmountLeft    float64   `json:"amount_left"`
	Symbol        string    `json:"symbol"`
	Type          string    `json:"type"`
	SellBuy       string    `json:"sellbuy"`
	Date          time.Time `json:"date"`
}

func CreateOrder(db *gorm.DB, userId uint, price, amount float64, symbol, orderType, sellbuy string) (Order, error) {
	order := Order{
		UserId:        userId,
		Price:         price,
		AmountInitial: amount,
		AmountLeft:    amount,
		Symbol:        symbol,
		Type:          orderType,
		SellBuy:       sellbuy,
		Date:          time.Now(),
	}
	if err := db.Create(&order).Error; err != nil {
		return Order{}, err
	}
	return order, nil
}
