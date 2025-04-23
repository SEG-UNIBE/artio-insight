package models

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Format used in the table `log_nip_tags`
type LogNipTags struct {
	Nip string
	Tag string
	TagStatus bool
	UpdateTime time.Time
}

// Transform the data retrieved from a `CurrentNipTags` object into a LogNipTags object, but only for one of the tags
func toLogNipTag(tag string, currentNipTags CurrentNipTags) (LogNipTags, error) {
	logNipTag := LogNipTags{
		Nip: currentNipTags.Nip,
		Tag: tag,
		TagStatus: false,
		UpdateTime: currentNipTags.LastUpdated,
	}

	switch tag {
		case "final":
			logNipTag.TagStatus = currentNipTags.Final
		case "draft":
			logNipTag.TagStatus = currentNipTags.Draft
		case "mandatory":
			logNipTag.TagStatus = currentNipTags.Mandatory
		case "optional":
			logNipTag.TagStatus = currentNipTags.Optional
		case "recommended":
			logNipTag.TagStatus = currentNipTags.Recommended
		case "unrecommended":
			logNipTag.TagStatus = currentNipTags.Unrecommended
		default:
			return LogNipTags{}, fmt.Errorf("unknown tag: %s", tag)
	}

	return logNipTag, nil
}

// Transform the data retrieved from a `CurrentNipTags` object into a list of LogNipTags objects, one for each tag
func toLogNipTagsList(currentNipTags CurrentNipTags) ([]LogNipTags, error) {
	finalTag, _ := toLogNipTag("final", currentNipTags)
	draftTag, _ := toLogNipTag("draft", currentNipTags)
	mandatoryTag, _ := toLogNipTag("mandatory", currentNipTags)
	optionalTag, _ := toLogNipTag("optional", currentNipTags)
	recommendedTag, _ := toLogNipTag("recommended", currentNipTags)
	unrecommendedTag, _ := toLogNipTag("unrecommended", currentNipTags)

	LogNipTagsList := []LogNipTags{finalTag, draftTag, mandatoryTag, optionalTag, recommendedTag, unrecommendedTag}

	return LogNipTagsList, nil
}

// Adds a new log entry in the database for a specific nip and tag
func addNewLogNipTag(db *gorm.DB, newTags CurrentNipTags, tag string) error {
	newDbLogTag, _ := toLogNipTag(tag, newTags)
	log.Info("updating \"", tag, "\" tag for NIP-", newTags.Nip)
	res := db.Create(newDbLogTag)

	if res.Error != nil {
		log.Error("could not add log entry for \"", tag, "\" of new NIP-", newTags.Nip, " : ", res.Error)
		return res.Error
	}

	return nil
}

// Compares the old tags with the current tags and adds the different ones to the log table in the database
func AddNewLogNipTags(db *gorm.DB, oldTags CurrentNipTags, newTags CurrentNipTags) bool {
	hasModified := false

	if oldTags.Final != newTags.Final {
		addNewLogNipTag(db, newTags, "final")
		hasModified = true
	}
	if oldTags.Draft != newTags.Draft {
		addNewLogNipTag(db, newTags, "draft")
		hasModified = true
	}
	if oldTags.Mandatory != newTags.Mandatory {
		addNewLogNipTag(db, newTags, "mandatory")
		hasModified = true
	}
	if oldTags.Optional != newTags.Optional {
		addNewLogNipTag(db, newTags, "optional")
		hasModified = true
	}
	if oldTags.Recommended != newTags.Recommended {
		addNewLogNipTag(db, newTags, "recommended")
		hasModified = true
	}
	if oldTags.Unrecommended != newTags.Unrecommended {
		addNewLogNipTag(db, newTags, "unrecommended")
		hasModified = true
	}

	return hasModified
}

// Add all the tags of a NIP to the log table in the database
func AddAllLogNipTags(db *gorm.DB, currNipTags CurrentNipTags) {
	logTags, logErr := toLogNipTagsList(currNipTags)
	if logErr != nil {
		log.Error(logErr)
	}

	var logResult *gorm.DB
	for _, logTag := range logTags {
		logResult = db.Create(&logTag)
		if logResult.Error != nil {
			log.Error("error while inserting log tag \"", logTag.Tag,"\" of new NIP-", currNipTags.Nip, " : ", logResult.Error)
		}
	}
}