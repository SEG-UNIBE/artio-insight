package main

import (
	m "nostr-protocol-nips_nip-tags/models"
	s "nostr-protocol-nips_nip-tags/pkg/nipTags"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Error("could not load .env file")
		return
	}

	err = m.InitDB()
	if err != nil {
		log.Error("could not initiate database connection: ", err)
		return
	}

	tags,_ := s.GetNipsTags()
	currentTime := time.Now()

	for nip, tag := range tags {
		nipTags,_ := m.ToCurrentNipTags(nip, currentTime, tag)

		result := m.DB.Create(&nipTags)
		if result.Error != nil {
			log.Error("error while inserting tags of NIP-", nip, " : ", result.Error)
		}
	}

	log.Info("data insertion finished")
}