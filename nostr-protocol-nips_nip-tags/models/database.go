package models

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	// Check if database already exists
	if DB != nil {
		return fmt.Errorf("DB already initialized")
	}

	// Retrieve database parameters
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	db_host := os.Getenv("DB_HOST")
	if db_user == "" || db_password == "" || db_name == "" || db_port == "" || db_host == "" {
		return fmt.Errorf("one or more required database environment variables are empty")
	}

	// Open database connection
	DSN := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Europe/Zurich", db_user, db_password, db_name, db_host, db_port)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	log.Info("Connected to database")

	return nil
}