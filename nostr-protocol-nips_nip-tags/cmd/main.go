package main

import (
	"fmt"
	m "nostr-protocol-nips_nip-tags/models"
	s "nostr-protocol-nips_nip-tags/pkg/nipTags"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	err = m.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	tags,_ := s.GetNipsTags()
	currentTime := time.Now()

	for nip, tag := range tags {
		nipTags,_ := m.ToCurrentNipTags(nip, currentTime, tag)

		result := m.DB.Create(&nipTags)
		if result.Error != nil {
			fmt.Println("Error inserting data:", result.Error)
		}
	}

	fmt.Println("\nData inserted successfully")
}