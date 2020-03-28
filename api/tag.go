package api

import (
	"strings"
	"time"
)

type Tag struct {
	Creator             int           `json:"creator"`
	Id                  int           `json:"id"`
	ImageId             interface{}   `json:"image_id"`
	Images              []TaggedImage `json:"images"`
	LastUpdated         time.Time     `json:"last_updated"`
	LastUpdater         int           `json:"last_updater"`
	LastUpdaterUsername string        `json:"last_updater_username"`
	//The name of the tag
	Name       string `json:"name"`
	Repository int    `json:"repository"`
	FullSize   int    `json:"full_size"`
	V2         bool   `json:"v2"`
}

type TaggedImage struct {
	Architecture string `json:"architecture"`
	Features     string `json:"features"`
	//Architecture variant - v7 for example
	Variant *string `json:"variant"`
	//Sha256 digest
	Digest     string  `json:"digest"`
	Os         string  `json:"os"`
	OsFeatures string  `json:"os_features"`
	OsVersion  *string `json:"os_version"`
	Size       uint32  `json:"size"`
}

func (t *Tag) String() string {
	return t.Name
}

type TagList []Tag

func (tags TagList) getByName(nm string, exactMatch bool) *Tag {
	for _, tag := range tags {
		if exactMatch {
			if tag.Name == nm {
				return &tag
			}
		} else {
			if strings.Contains(tag.Name, nm) {
				return &tag
			}
		}
	}
	return nil
}
