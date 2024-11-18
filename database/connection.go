package database

import (
	"jwt/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=" + os.Getenv("DBUSER") + " password=" + os.Getenv("DBPASS") + " dbname=" + os.Getenv("DBNAME") + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// dsn := "root:" + os.Getenv("DBPASS") + "@tcp(127.0.0.1:3306)/" + os.Getenv("DBNAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = connection

	// connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Order{})
}
