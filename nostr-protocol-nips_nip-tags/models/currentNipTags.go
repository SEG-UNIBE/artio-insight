package models

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Format used in the table `current_nip_tags`
type CurrentNipTags struct {
	Nip string
	Final bool
	Draft bool
	Mandatory bool
	Optional bool
	Recommended bool
	Unrecommended bool
	LastUpdated time.Time
}

// Transforms the data retrieved from the github repo into a CurrentNipTags object
func ToCurrentNipTags(nip string, lastUpdated time.Time, tags []string) (CurrentNipTags, error) {
	cnt := CurrentNipTags{
		Nip: nip,
		Final: false,
		Draft: false,
		Mandatory: false,
		Optional: false,
		Recommended: false,
		Unrecommended: false,
		LastUpdated: lastUpdated,
	}
	var err error = nil

	for _, t := range tags {
		switch t {
			case "final":
				cnt.Final = true
			case "draft":
				cnt.Draft = true
			case "mandatory":
				cnt.Mandatory = true
			case "optional":
				cnt.Optional = true
			case "recommended":
				cnt.Recommended = true
			case "unrecommended":
				cnt.Unrecommended = true
			default:
				err = fmt.Errorf("unknown tag: %s", t)
		}
	}

	log.Error(err)
	return cnt, err
}

// Add the current nip tags to the database
func AddCurrentNipTags(db *gorm.DB, nipTags CurrentNipTags) {
	log.Info("creating nip tags for NIP-", nipTags.Nip)

	res := db.Create(&nipTags)
	if res.Error != nil {
		log.Error("error while inserting tags of new NIP-", nipTags.Nip, " : ", res.Error)
	}
}

// Update the tags and time of the given nip in the database
func UpdateCurrentNipTags(db *gorm.DB, nipTags CurrentNipTags) {
	res := db.Where("nip = ?", nipTags.Nip).Save(&nipTags)
	if res.Error != nil {
		log.Error("error while updating tags of NIP-", nipTags.Nip, " : ", res.Error)
	}
}