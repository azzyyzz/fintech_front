package models

import (
	"gorm.io/gorm"
)

func Deposit(db *gorm.DB, userId uint, symbol string, amount float64) error {
	var wallet Wallet
	if err := db.Where("user_id = ? AND symbol = ?", userId, symbol).First(&wallet).Error; err != nil {
		return err
	}
	wallet.Balance += amount
	return db.Save(&wallet).Error
}
