package api

import (
	"regexp"
	"strings"
	"time"
)

type UserRepository struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type RepositoryPermissions struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
	Admin bool `json:"admin"`
}

type Repository struct {
	User            string                 `json:"user"`
	Name            string                 `json:"name"`
	Namespace       string                 `json:"namespace"`
	RepositoryType  string                 `json:"repository_type"`
	Status          int                    `json:"status"`
	Description     string                 `json:"description"`
	IsPrivate       bool                   `json:"is_private"`
	IsAutomated     bool                   `json:"is_automated"`
	CanEdit         bool                   `json:"can_edit"`
	StarCount       int                    `json:"start_count"`
	PullCount       int                    `json:"pull_count"`
	LastUpdated     *time.Time             `json:"last_updated"`
	IsMigrated      bool                   `json:"is_migrated"`
	HasStarred      bool                   `json:"has_starred"`
	FullDescription string                 `json:"full_description"`
	Affiliation     string                 `json:"affiliation"`
	Permissions     *RepositoryPermissions `json:"permissions"`
}

//GetGitRepoLinks
func (r *Repository) GetGitRepoLinks() []string {
	d := r.FullDescription
	var rx, err = regexp.Compile("https?://(www\\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}([-a-zA-Z0-9()@:%_+.~#?&/=]*)")
	if err != nil {
		return nil
	}
	matches := rx.FindAllString(d, -1)
	var output []string
	var repoSites = []string{"github.com", "butbucket.org"}
	for _, m := range matches {
		m = strings.TrimRight(m, ").")
		if strings.ContainsRune(m, '#') {
			continue
		}
		isRepoSite := false
		for _, reposite := range repoSites {
			if strings.Contains(m, reposite) {
				isRepoSite = true
				break
			}
		}
		if isRepoSite {
			slashparts := strings.SplitN(m, "/", 7)
			//We dont need long urls.
			//presumed url is https://github.com/owner/repo
			if len(slashparts) != 5 {
				continue
			}
			output = append(output, m)
		}
	}
	return output
}
