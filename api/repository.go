package api

import (
	"fmt"
	"regexp"
	"sort"
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

type NamedLink struct {
	Name string
	Link string
}

func (n *NamedLink) String() string {
	return fmt.Sprintf("%s - %s ", n.Name, n.Link)
}

func (r *Repository) GetLinks() []string {
	return getLinks(r.FullDescription)
}

func (r *Repository) GetMarkdownLinks() []NamedLink {
	return getMarkdownLinks(r.FullDescription)
}

func getLinks(str string) []string {
	var rx, err = regexp.Compile("https?://(www\\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}([-a-zA-Z0-9()@:%_+.~#?&/=]*)")
	if err != nil {
		return nil
	}
	matches := rx.FindAllString(str, -1)
	return matches
}

func getMarkdownLinks(strMarkdown string) []NamedLink {
	rx, _ := regexp.Compile("\\[([^]]+)\\]\\((https?://\\S+)\\)")
	matches := rx.FindAllStringSubmatch(strMarkdown, -1)
	var output []NamedLink
	for _, m := range matches {
		output = append(output, NamedLink{m[1], m[2]})
	}
	return output
}

func stringIsMarkdown(str string) bool {
	hasLinks, err := regexp.MatchString("\\[[^]]+\\]\\((https?://\\S+)\\)", str)
	if err != nil {
		return false
	}
	return hasLinks
}

func isRepoSite(link string) bool {
	slashparts := strings.SplitN(link, "/", 7)
	//We dont need long urls.
	//presumed url is https://github.com/owner/repo
	if len(slashparts) != 5 {
		return false
	}
	var repoSites = []string{"github.com", "butbucket.org"}
	for _, reposite := range repoSites {
		if strings.Contains(link, reposite) {
			return true
		}
	}
	return false
}

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func getMostCommonUrl(urls []string, slashCount int) string {
	freqs := make(map[string]int)
	for _, url := range urls {
		parts := strings.SplitN(url, "/", slashCount+1)
		trimmedUrl := strings.Join(parts[0:slashCount], "/")
		freqs[trimmedUrl] += 1
	}
	i := 0
	s := make(PairList, len(freqs))
	for k, v := range freqs {
		s[i] = Pair{k, v}
	}
	sort.Sort(sort.Reverse(s))
	return s[0].Key
}

func (r *Repository) GetGitRepo() string {
	validLinks := r.GetGitRepoLinks()
	url := getMostCommonUrl(validLinks, 5)
	return url
}

//GetGitRepoLinks gets all the links that are rleated to a github/bitbucket repo
func (r *Repository) GetGitRepoLinks() []string {
	d := r.FullDescription
	var matches []string
	if stringIsMarkdown(d) {
		links := getMarkdownLinks(d)
		for _, l := range links {
			matches = append(matches, l.Link)
		}
	} else {
		matches = getLinks(d)
	}
	var validLinks []string
	for _, l := range matches {
		if strings.ContainsRune(l, '#') {
			continue
		}
		isDockerfile := strings.HasSuffix(l, "/Dockerfile")
		isIssues := strings.HasSuffix(l, "/issues")
		if isDockerfile || isIssues || isRepoSite(l) {
			validLinks = append(validLinks, l)
		}
	}
	return validLinks
}
