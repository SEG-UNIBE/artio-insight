package models

import (
	"fmt"
	"time"
)

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

func GetCurrentNipTags(nip string, lastUpdated time.Time, tags []string) (CurrentNipTags, error) {
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

	return cnt, err
}