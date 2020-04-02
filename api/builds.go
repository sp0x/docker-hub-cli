package api

import (
	"errors"
	"fmt"
	"net/url"
)

type BuildSource struct {
	AutoTests   string `json:"autotests"`
	BuildInFarm bool   `json:"build_in_farm"`
	Channel     string `json:"channel"`
	Image       string `json:"image"`
	Owner       string `json:"owner"`
	Provider    string `json:"provider"`
	RepoLinks   bool   `json:"repo_links"`
	Repository  string `json:"repository"`
	ResourceUri string `json:"resource_uri"`
	State       string `json:"state"`
	Uuid        string `json:"uuid"`
}

func (b *BuildSource) GetSourceUrl() (*url.URL, error) {
	rurl := url.URL{}
	rurl.Scheme = "https"
	if b.Provider == "Github" {
		rurl.Host = "github.com"
		rurl.Path = fmt.Sprintf("%s/%s", b.Owner, b.Repository)
	} else {
		return nil, errors.New("build provider not supported")
	}
	return &rurl, nil
}
