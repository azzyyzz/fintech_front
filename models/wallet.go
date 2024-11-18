package models

import (
	"gorm.io/gorm"
)

type Wallet struct {
	Id      uint    `json:"id" gorm:"primaryKey"`
	UserId  uint    `json:"user_id"`
	Symbol  string  `json:"symbol" gorm:"type:varchar(10)"`
	Balance float64 `json:"balance"`
}

func CreateWalletsForUser(db *gorm.DB, userId uint) {
	symbols := []string{"USDT", "BTC", "ETH"}
	for _, symbol := range symbols {
		wallet := Wallet{
			UserId:  userId,
			Symbol:  symbol,
			Balance: 0, // Initial balance is set to 0
		}
		db.Create(&wallet)
	}
}
