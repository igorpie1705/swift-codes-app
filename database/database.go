package database

import (
	"log"

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
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err := db.AutoMigrate(&models.SwiftCode{}); err != nil {
		log.Fatalf("Error during database migration: %v", err)
	}

	log.Println("Database connection established and migration completed successfully.")
	return db
}

func GetDB() *gorm.DB {
	return db
}

func SetDB(database *gorm.DB) {
	db = database
}
