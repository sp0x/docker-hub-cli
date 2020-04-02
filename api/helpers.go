package api

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type NamedLink struct {
	Name string
	Link string
}

func (n *NamedLink) String() string {
	return fmt.Sprintf("%s - %s ", n.Name, n.Link)
}

func getMostCommonUrl(urls []string, slashCount int) string {
	freqs := make(map[string]int)
	for _, urlx := range urls {
		parts := strings.SplitN(urlx, "/", slashCount+1)
		trimmedUrl := strings.Join(parts[0:slashCount], "/")
		freqs[trimmedUrl] += 1
	}
	i := 0
	s := make(PairList, len(freqs))
	for k, v := range freqs {
		s[i] = Pair{k, v}
		i += 1
	}
	sort.Sort(sort.Reverse(s))
	if len(s) == 0 {
		return ""
	}
	return s[0].Key
}

func getLinks(str string) []string {
	var rx, err = regexp.Compile("https?://(www\\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}([-a-zA-Z0-9@:%_+.~#?&/=]*)")
	if err != nil {
		return nil
	}
	matches := rx.FindAllString(str, -1)
	for ix, m := range matches {
		matches[ix] = strings.TrimRight(m, ")")
	}
	return matches
}

func filterLinks(links []NamedLink, filter func(link NamedLink) bool) []NamedLink {
	var output []NamedLink
	for _, m := range links {
		if filter != nil {
			if !filter(m) {
				continue
			}
		}
		output = append(output, m)
	}
	return output
}

//Gets all the markdown links from a given text.
func getMarkdownLinks(strMarkdown string, filter func(link NamedLink) bool) []NamedLink {
	rx, _ := regexp.Compile("\\[([^]]+)\\]\\((https?://\\S+)\\)")
	matches := rx.FindAllStringSubmatch(strMarkdown, -1)
	var output []NamedLink
	for _, m := range matches {
		nlink := NamedLink{m[1], m[2]}
		if filter != nil {
			if !filter(nlink) {
				continue
			}
		}
		output = append(output, nlink)
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
