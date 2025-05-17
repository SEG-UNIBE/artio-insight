package models

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const SQL_FILE_PATH = "init/init.sql"

var DB *gorm.DB

func InitDB(isDocker bool) error {
	// Check if database already exists
	if DB != nil {
		return fmt.Errorf("database already initialized")
	}

	// Retrieve database parameters
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	if isDocker {
		dbPort = os.Getenv("DOCKER_DB_PORT")
		dbHost = os.Getenv("DOCKER_DB_HOST")
	}
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
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	  })
	if err != nil {
		return err
	}

	DB = db
	log.Info("connected to database")

	// Add the schema to the database
	err = DB.Exec("BEGIN;").Error
	if err != nil {
		return err
	}

	sqlContent, readErr := os.ReadFile(SQL_FILE_PATH)
	if readErr != nil {
		return readErr
	}

	err = DB.Exec(string(sqlContent)).Error
	if err != nil {
		DB.Exec("ROLLBACK;")
		return err
	}

	err = DB.Exec("COMMIT;").Error
	if err != nil {
		return err
	}
	
	return nil
}