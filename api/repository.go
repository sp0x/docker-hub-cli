package api

import (
	"errors"
	"fmt"
	"net/url"
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
	isMd            *bool
	tags            TagList
	markdownLinks   []NamedLink
}

//GetLinks Gets the links from the full description
func (r *Repository) GetLinks() []string {
	return getLinks(r.FullDescription)
}

//GetMarkdownLinks gets all the markdown links from the repository's full description
func (r *Repository) GetMarkdownLinks() []NamedLink {
	if r.markdownLinks != nil {
		return r.markdownLinks
	}
	r.markdownLinks = getMarkdownLinks(r.FullDescription, nil)
	return r.markdownLinks
}

//GetTaggedDockerfile gets a link to the dockerfile for a given tag.
//This works only with repositories that have Markdown descriptions, since the link names are used too in order to figure out the tag
func (r *Repository) GetTaggedDockerfile(dapi *DockerApi, tagName string, exactTagMatch bool) (string, error) {
	if !r.IsMarkdowned() {
		return "", errors.New("description is not in markdown")
	}
	//Check if there's a tag with that name
	tags, err := dapi.GetTagsFromRepo(r, 0, 0)
	if err != nil {
		return "", err
	}
	tag := tags.getByName(tagName, exactTagMatch)
	if tag == nil {
		return "", errors.New("tag not found")
	}
	mlinks := filterLinks(r.GetMarkdownLinks(), func(l NamedLink) bool {
		isDockerfile := strings.HasSuffix(l.Link, "/Dockerfile")
		//If it's a dockerfile link or the link's name contains the tag name enclosed in ``
		return isDockerfile && strings.Contains(l.Name, fmt.Sprintf("`%s`", tagName))
	})
	if mlinks == nil || len(mlinks) == 0 {
		return "", errors.New("the repository has no links in it's description")
	}
	return mlinks[0].Link, nil
}

//IsMarkdowned checks if there's any signs that the full description of the repository is in markdown.
func (r *Repository) IsMarkdowned() bool {
	if r.isMd != nil {
		return *r.isMd
	}
	res := stringIsMarkdown(r.FullDescription)
	r.isMd = &res
	return res
}

//GetTaggedRepositoryDirectory gets the directory path for a given tag, from a git repository for this repository
func (r *Repository) GetTaggedRepositoryDirectory(dapi *DockerApi, tagName string, exactTagMatch bool) (string, error) {
	dockerfile, err := r.GetTaggedDockerfile(dapi, tagName, exactTagMatch)
	if err != nil {
		return "", err
	}
	//Most times this is the directory of the dockerfile
	durl, err := url.Parse(dockerfile)
	if err != nil {
		return "", err
	}
	dflen := len("/Dockerfile")
	if durl.Path == "" || len(durl.Path) <= dflen {
		return "", errors.New("bad dockerfile url")
	}
	pth := string(durl.Path[0 : len(durl.Path)-dflen])
	//Github uses urls like /nginxinc/docker-nginx/blob/5c15613519a26c6adc244c24f814a95c786cfbc3/mainline/alpine/Dockerfile
	if strings.Contains(pth, "/blob/") && durl.Host == "github.com" {
		pathParts := strings.Split(pth, "/")
		if len(pathParts) > 5 && pathParts[3] == "blob" {
			pth = strings.Join(pathParts[:3], "/") + "/" + strings.Join(pathParts[5:], "/")
		}
	}
	return pth, nil
}

//GetGitRepo gets a link to the git repository for this repository
func (r *Repository) GetGitRepo() string {
	validLinks := r.GetGitRepoLinks()
	if validLinks == nil {
		return ""
	}
	urlx := getMostCommonUrl(validLinks, 5)
	return urlx
}

//GetGitRepoLinks gets all the links that are rleated to a github/bitbucket repo
func (r *Repository) GetGitRepoLinks() []string {
	d := r.FullDescription
	var matches []string
	matches = getLinks(d)
	//if r.IsMarkdowned() {
	//	links := r.GetMarkdownLinks()
	//	for _, l := range links {
	//		matches = append(matches, l.Link)
	//	}
	//} else {
	//	matches = getLinks(d)
	//}
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
