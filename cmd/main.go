package main

import (
	m "artio-insight/models"
	s "artio-insight/pkg/scripts"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	g "gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")

	DSN := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=%s sslmode=disable TimeZone=Europe/Zurich", db_user, db_password, db_name, db_port)

	db, err := g.Open(postgres.New(postgres.Config{
		DSN: DSN,
		PreferSimpleProtocol: true,
	}), &g.Config{})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	fmt.Println("Connected to database")

	tags,_ := s.GetNipsTags()
	currentTime := time.Now()

	for nip, tag := range tags {
		nipTags,_ := m.GetCurrentNipTags(nip, currentTime, tag)

		result := db.Create(&nipTags)
		if result.Error != nil {
			fmt.Println("Error inserting data:", result.Error)
		}
	}

	fmt.Println("\nData inserted successfully")
}