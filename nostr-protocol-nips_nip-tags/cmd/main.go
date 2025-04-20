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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Error("could not load .env file")
		return
	}

	err = db.InitDB()
	if err != nil {
		log.Error("could not initiate database connection: ", err)
		return
	}

	// TODO : move out of main
	rawTags, _ := s.GetNipsTags()
	currentTime := time.Now().UTC()
	hasModified := false

	for nip, tags := range rawTags {
		nipTags, currErr := m.ToCurrentNipTags(nip, currentTime, tags)
		if currErr != nil {
			log.Error(currErr)
		}

		var dbNipTags m.CurrentNipTags
		dbErr := db.DB.Table("current_nip_tags").Take(&dbNipTags, "nip = ?", nip)
		if errors.Is(dbErr.Error, gorm.ErrRecordNotFound) {
			log.Info("creating nip tags for NIP-", nip)
			
			result := db.DB.Create(&nipTags)
			if result.Error != nil {
				log.Error("error while inserting tags of new NIP-", nip, " : ", result.Error)
			}
			
			logTags, logErr := m.ToLogNipTagsList(nipTags)
			if logErr != nil {
				log.Error(logErr)
			}
			var logResult *gorm.DB
			for _, logTag := range logTags {
				logResult = db.DB.Create(&logTag)
				if logResult.Error != nil {
					log.Error("error while inserting log tag \"", logTag.Tag,"\" of new NIP-", nip, " : ", logResult.Error)
				}
			}
		} else if dbErr.Error != nil {
			log.Error("could not check if nip ", nip, " exists in database: ", dbErr.Error)
		} else  {
			if dbNipTags.Final != nipTags.Final {
				newDbLogTag, _ := m.ToLogNipTags("final", nipTags)
				log.Info("updating \"final\" tag for NIP-", nip)
				hasModified = true
				db.DB.Create(newDbLogTag)
			}
			if dbNipTags.Draft != nipTags.Draft {
				newDbLogTag, _ := m.ToLogNipTags("draft", nipTags)
				log.Info("updating \"draft\" tag for NIP-", nip)
				hasModified = true
				db.DB.Create(newDbLogTag)
			}
			if dbNipTags.Mandatory != nipTags.Mandatory {
				newDbLogTag, _ := m.ToLogNipTags("mandatory", nipTags)
				log.Info("updating \"mandatory\" tag for NIP-", nip)
				hasModified = true
				db.DB.Create(newDbLogTag)
			}
			if dbNipTags.Optional != nipTags.Optional {
				newDbLogTag, _ := m.ToLogNipTags("optional", nipTags)
				log.Info("updating \"optional\" tag for NIP-", nip)
				hasModified = true
				db.DB.Create(newDbLogTag)
			}
			if dbNipTags.Recommended != nipTags.Recommended {
				newDbLogTag, _ := m.ToLogNipTags("recommended", nipTags)
				log.Info("updating \"recommended\" tag for NIP-", nip)
				hasModified = true
				db.DB.Create(newDbLogTag)
			}
			if dbNipTags.Unrecommended != nipTags.Unrecommended {
				newDbLogTag, _ := m.ToLogNipTags("unrecommended", nipTags)
				log.Info("updating \"unrecommended\" tag for NIP-", nip)
				hasModified = true
				db.DB.Create(newDbLogTag)
			}

			res := db.DB.Where("nip = ?", nip).Save(&nipTags)
			if res.Error != nil {
				log.Error("error while updating tags of NIP-", nip, " : ", res.Error)
			}
		}
	}

	if hasModified {
		log.Info("tag insertion finished with modifications")
	} else {
		log.Info("tag insertion finished without modifications")
	}
}