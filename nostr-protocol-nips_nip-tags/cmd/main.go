package main

import (
	"errors"
	m "nostr-protocol-nips_nip-tags/models"
	s "nostr-protocol-nips_nip-tags/pkg/nipTags"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

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

		curr := m.DB.First(&m.CurrentNipTags{}, nip)
		if errors.Is(curr.Error, gorm.ErrRecordNotFound) {
			log.Info("new nip found : NIP-", nip)
			result := m.DB.Create(&nipTags)
			if result.Error != nil {
				log.Error("error while inserting tags of new NIP-", nip, " : ", result.Error)
			}
			continue
		} else if curr.Error != nil {
			log.Error("could not check if nip ", nip, " exists in database: ", curr.Error)
			continue
		}
	}

	log.Info("data insertion finished")
}