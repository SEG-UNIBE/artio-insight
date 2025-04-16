package models

import (
	"fmt"
	"time"
)

// Format used in the table `log_nip_tags`
type LogNipTags struct {
	Nip string
	Tag string
	TagStatus bool
	UpdateTime time.Time
}

// Transform the data retrieved from a `CurrentNipTags` object into a LogNipTags object, but only for one of the tags
func ToLogNipTags(tag string, currentNipTags CurrentNipTags) (LogNipTags, error) {
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
func ToLogNipTagsList(currentNipTags CurrentNipTags) ([]LogNipTags, error) {
	finalTag, _ := ToLogNipTags("final", currentNipTags)
	draftTag, _ := ToLogNipTags("draft", currentNipTags)
	mandatoryTag, _ := ToLogNipTags("mandatory", currentNipTags)
	optionalTag, _ := ToLogNipTags("optional", currentNipTags)
	recommendedTag, _ := ToLogNipTags("recommended", currentNipTags)
	unrecommendedTag, _ := ToLogNipTags("unrecommended", currentNipTags)

	LogNipTagsList := []LogNipTags{finalTag, draftTag, mandatoryTag, optionalTag, recommendedTag, unrecommendedTag}

	return LogNipTagsList, nil
}
