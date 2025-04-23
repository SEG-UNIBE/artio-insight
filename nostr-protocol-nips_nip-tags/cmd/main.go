package main

import (
	"errors"
	m "nostr-protocol-nips_nip-tags/models"
	db "nostr-protocol-nips_nip-tags/pkg/database"
	s "nostr-protocol-nips_nip-tags/pkg/nipTags"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	// Get Environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Error("could not load .env file")
		return
	}

	// Get the database connection
	err = db.InitDB()
	if err != nil {
		log.Error("could not initiate database connection: ", err)
		return
	}

	// Retrieve the NIP Tags from Github
	rawTags, _ := s.GetNipsTags()

	// Add the tags to the database
	currentTime := time.Now().UTC()
	hasModified := false
	for nip, tags := range rawTags {
		newNipTags, _ := m.ToCurrentNipTags(nip, currentTime, tags)

		var oldNipTags m.CurrentNipTags
		dbErr := db.DB.Table("current_nip_tags").Take(&oldNipTags, "nip = ?", nip)

		// If a new tag, add everything
		if errors.Is(dbErr.Error, gorm.ErrRecordNotFound) {
			log.Info("creating nip tags for NIP-", nip)
			hasModified = true
			
			m.AddCurrentNipTags(db.DB, newNipTags)
			m.AddAllLogNipTags(db.DB, newNipTags)
		// If error while checking nip, do nothing
		} else if dbErr.Error != nil {
			log.Error("could not check if nip ", nip, " exists in database: ", dbErr.Error)
		// If nip exists, only add the difference from previous pooling
		} else  {
			hasModified = m.AddNewLogNipTags(db.DB, newNipTags, oldNipTags)
			m.UpdateCurrentNipTags(db.DB, newNipTags)
		}
	}

	if hasModified {
		log.Info("tag insertion finished with modifications")
	} else {
		log.Info("tag insertion finished without modifications")
	}
}