package main

import (
	m "artio-insight/models"
	s "artio-insight/pkg/scripts"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	g "gorm.io/gorm"
)

func main() {
	db, err := g.Open(postgres.New(postgres.Config{
		DSN: "user=postgres password=admin dbname=artio_db host=localhost port=5432 sslmode=disable TimeZone=Europe/Zurich",
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