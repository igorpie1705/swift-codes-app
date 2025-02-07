package database

import (
	"github.com/igorpie1705/swift-codes-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() *gorm.DB {
	dsn := "host=db user=user password=password dbname=swiftcodes port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&models.SwiftCode{})

	return db
}

func GetDB() *gorm.DB {
	return db
}