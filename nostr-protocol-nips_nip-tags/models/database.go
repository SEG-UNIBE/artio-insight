package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error{
	if DB != nil {
		return fmt.Errorf("DB already initialized")
	}

	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")

	DSN := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=%s sslmode=disable TimeZone=Europe/Zurich", db_user, db_password, db_name, db_port)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}

	DB = db

	fmt.Println("Connected to database")

	return nil
}