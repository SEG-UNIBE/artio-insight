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
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbSslmode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")
	if dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbHost == "" || dbSslmode == "" || dbTimezone == "" {
		return fmt.Errorf("one or more required database environment variables are empty")
	}

	// Open database connection
	DSN := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s", dbUser, dbPassword, dbName, dbHost, dbPort, dbSslmode, dbTimezone)
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