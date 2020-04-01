package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/docker-hub-cli/requests"
	"net/http"
	"net/http/cookiejar"
	"path"
	"strings"
	"time"
)

func NewApi() *DockerApi {
	d := &DockerApi{}
	version := "2"
	transport := &http.Transport{
		DisableCompression: false,
	}
	//Cookies are needed for authentication
	cookies, _ := cookiejar.New(nil)
	d.client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
		//Jar:       cookies, //Commented because this causes CSRF issues if enabled
	}
	//To keep a session
	d.cookieJar = cookies
	d.apiVersion = version
	d.routeBase = fmt.Sprintf("https://hub.docker.com/v%s", version)
	return d
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

type DockerApi struct {
	client     *http.Client
	apiVersion string
	routeBase  string
	token      string
	cookieJar  *cookiejar.Jar
	username   string
}

func (d *DockerApi) getRoute(p string) string {
	return joinURL(d.routeBase, p)
}

func (d *DockerApi) SetRepositoryDescription(username, name string, descShort, descLong string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s", username, name)) + "/"
	data := map[string]string{}
	if descLong != "" {
		data["full_description"] = descLong
	}
	if descShort != "" {
		data["description"] = descShort
	}
	r, err := requests.Patch(d.client, pth, data, d.token)
	if err != nil {
		if r != nil {
			return fmt.Errorf(parseError(r))
		} else {
			return err
		}
	}
	return nil
}

func (d *DockerApi) SetRepositoryPrivacy(username, name string, isPrivate bool) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/privacy", username, name))
	r, err := requests.Post(d.client, pth, map[string]interface{}{
		"is_private": isPrivate,
	}, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(string(r))
	return nil
}

//GetBuildSettings gets the build settings for a repository
func (d *DockerApi) GetBuildSettings(username string, name string) (string, error) {
	username = strings.ToLower(username)
	settingsPath := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild", username, name))
	r, err := requests.Get(d.client, settingsPath, d.token)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return string(r), nil
}

//GetBuildDetails Gets the details for a given build of a repository.
func (d *DockerApi) GetBuildDetails(username, name, code string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if code == "" {
		return fmt.Errorf("no build code given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	settingsPath := d.getRoute(fmt.Sprintf("repositories/%s/%s/buildhistory/%s", username, name, code))
	r, err := requests.Get(d.client, settingsPath, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(string(r))
	return nil
}

//MyRepositories gets the repositories of the currently logged in user.
func (d *DockerApi) MyRepositories() ([]UserRepository, error) {
	if d.username == "" {
		return nil, fmt.Errorf("user not authenticated")
	}
	return d.Repositories(d.username)
}

//Repositories gets the repositories of an user
func (d *DockerApi) Repositories(username string) ([]UserRepository, error) {
	if username == "" {
		return nil, fmt.Errorf("no user given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("users/%s/repositories", username))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	var repositories []UserRepository
	err = json.Unmarshal(r, &repositories)
	if err != nil {
		return nil, err
	}
	return repositories, nil
}

//GetRepositoriesStarred Gets the starred repositories for a user.
func (d *DockerApi) GetRepositoriesStarred(username string, page, pageSize int) ([]UserRepository, error) {
	if username == "" {
		return nil, fmt.Errorf("no user given")
	}
	username = strings.ToLower(username)
	if page == 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}

	pth := d.getRoute(fmt.Sprintf("users/%s/repositories/starred?page_size=%v&page=%v", username, pageSize, page))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	var search SearchResult
	err = json.Unmarshal(r, &search)
	if err != nil {
		return nil, err
	}
	if search.Count == 0 {
		return nil, nil
	}
	var results []UserRepository
	err = json.Unmarshal(search.Results, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

//GetBuildTriggerHistory Gets the build trigger history for a given repository.
func (d *DockerApi) GetBuildTriggerHistory(username, name string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no repo name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/buildtrigger/history", username, name))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return err
	}
	log.Print(string(r))
	return nil
}

//GetBuildTrigger Gets the build trigger for a given repository.
func (d *DockerApi) GetBuildTrigger(username, name string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/buildtrigger", username, name))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(string(r))
	return nil
}

func (d *DockerApi) SaveBuildTag(username, name string, id string, tagName, dockerfileLocation, sourceType, sourceName string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild/tags/%s", username, name, id))
	if dockerfileLocation == "" {
		dockerfileLocation = "/"
	}
	if sourceType == "" {
		sourceType = "Branch"
	}
	if sourceName == "" {
		sourceName = "master"
	}
	data := map[string]string{
		"id":                  id,
		"name":                tagName,
		"dockerfile_location": dockerfileLocation,
		"source_type":         sourceType,
		"source_name":         sourceName,
	}
	r, err := requests.Put(d.client, pth, data, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(string(r))
	return nil
}

//GetComments gets the comments for an image, default  page size is 100, pages start from 1
func (d *DockerApi) GetComments(username, name string, pageSize int, page int) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/comments?page_size=%v&page=%v", username, name, pageSize, page))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

//TODO Gets the tags for a repository
func (d *DockerApi) GetTags(username, name string, pageSize, page int) (TagList, error) {
	if username != "" && name == "" {
		name = username
		username = "library "
	}
	if username == "" || username == "_" {
		username = "library"
	}
	username = strings.ToLower(username)
	if pageSize == 0 {
		pageSize = 100
	}
	if page < 1 {
		page = 1
	}
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/tags?page_size=%v&page=%v", username, name, pageSize, page))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	var searchRes SearchResult
	err = json.Unmarshal(r, &searchRes)
	if err != nil {
		return nil, err
	}
	var tags []Tag
	err = json.Unmarshal(searchRes.Results, &tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

//GetMyRepository gets details about a user owned repository
func (d *DockerApi) GetMyRepository(name string) (*Repository, error) {
	if d.username == "" {
		return nil, fmt.Errorf("user not authenticated")
	}
	return d.GetRepository(d.GetUsername(), name)
}

//GetRepository gets details about a repository
func (d *DockerApi) GetRepository(username, name string) (*Repository, error) {
	var repo Repository
	if username != "" && name == "" {
		name = username
		username = "library"
	}
	if username == "_" || username == "" {
		username = "library"
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s", username, name))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

//CreateBuildLink Creates a build link for a given repository to the given repository.
func (d *DockerApi) CreateBuildLink(username, name, toRepo string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if toRepo == "" {
		return fmt.Errorf("no target repo given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	if !strings.Contains(toRepo, "/") {
		toRepo = fmt.Sprintf("library/%s", toRepo)
	}
	if string(toRepo[0:2]) == "_/" {
		toRepo = fmt.Sprintf("library/%s", string(toRepo[2:]))
	}
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/links", username, name))
	data := map[string]string{
		"to_repo": toRepo,
	}
	_, err := requests.Post(d.client, pth, data, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	return nil
}

//CreateBuildTag Creates a build tag for a given repository.
func (d *DockerApi) CreateBuildTag(username, name, tagname, dockerFileLocation, sourceType, sourceName string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if tagname == "" {
		return fmt.Errorf("no tagname given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	if dockerFileLocation == "" {
		dockerFileLocation = "/"
	}
	if sourceType == "" {
		sourceType = "Branch"
	}
	if sourceName == "" {
		sourceName = "master"
	}

	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild/tags", username, name))
	data := map[string]interface{}{
		"isNew":               true,
		"namespace":           username,
		"repoName":            name,
		"name":                tagname,
		"dockerfile_location": dockerFileLocation,
		"source_type":         sourceType,
		"source_name":         sourceName,
	}
	r, err := requests.Post(d.client, pth, data, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

//CreateAutomatedBuild - Creates an automated build.
func (d *DockerApi) CreateAutomatedBuild(username, name string, details map[string]string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild/", username, name))
	data := map[string]interface{}{
		"name":                name,
		"namespace":           username,
		"active":              true,
		"dockerhub_repo_name": fmt.Sprintf("%s/%s", username, name),
		"is_private":          false,
	}
	//Fill details
	for k, v := range details {
		data[k] = v
	}
	r, err := requests.Post(d.client, pth, data, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

func (d *DockerApi) CreateOwnRepository(name string, isPrivate bool, desc, fullDesc string) (*Repository, error) {
	if d.username == "" {
		return nil, fmt.Errorf("user not authenticated")
	}
	return d.CreateRepository(d.username, name, isPrivate, desc, fullDesc)
}

//CreateRepository creates a repository
func (d *DockerApi) CreateRepository(username, name string, isPrivate bool, desc, fullDesc string) (*Repository, error) {
	if username == "" {
		return nil, fmt.Errorf("no user given")
	}
	if name == "" {
		return nil, fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute("repositories") + "/"
	data := map[string]interface{}{
		"name":             name,
		"namespace":        username,
		"is_private":       isPrivate,
		"description":      desc,
		"full_description": fullDesc,
	}
	r, err := requests.Post(d.client, pth, data, d.token)
	if err != nil {
		if r != nil {
			return nil, fmt.Errorf(parseError(r))
		} else {
			return nil, err
		}
	}
	var repo Repository
	err = json.Unmarshal(r, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

//DeleteBuildLink Deletes a build link for a given repository.
func (d *DockerApi) DeleteBuildLink(username, name, id string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/links/%s", username, name, id))
	r, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

//DeleteBuildTag Deletes a build tag for a given repository.
func (d *DockerApi) DeleteBuildTag(username, name, id string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if id == "" {
		return fmt.Errorf("no tag id given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild/tags/%s", username, name, id))
	r, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

//DeleteCollaborator - Deletes a build tag for a given repository.
func (d *DockerApi) DeleteCollaborator(username, name, collaborator string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if collaborator == "" {
		return fmt.Errorf("no collaborator username given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/collaborators/%s", username, name, collaborator))
	r, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

func (d *DockerApi) DeleteOwnRepository(name string) error {
	if d.username == "" {
		return fmt.Errorf("user not authenticated")
	}
	return d.DeleteRepository(d.username, name)
}

//DeleteRepository Deletes a repository.
func (d *DockerApi) DeleteRepository(username, name string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s", username, name)) + "/"
	_, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		return err
	}
	return nil
}

//DeleteTag - Deletes a tag for the given username and repository.
func (d *DockerApi) DeleteTag(username, name, tag string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if tag == "" {
		return fmt.Errorf("no tag name given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/tags/%s", username, name, tag))
	r, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

//GetRegistrySettings gets the settings for the current logged in user containing information about the number of private repositories used/available.
func (d *DockerApi) GetRegistrySettings(username string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("users/%s/registry-settings", username))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
}

//TODO Creates a build tag for a given repository.
func (d *DockerApi) TriggerBuild(username, name string, dockerfileLocation, sourceType, sourceName string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild/trigger-build", username, name))
	data := map[string]string{
		"dockerfile_location": dockerfileLocation,
		"source_type":         sourceType,
		"source_name":         sourceName,
	}
	r, err := requests.Post(d.client, pth, data, d.token)
	if err != nil {
		return nil
	}
	log.Print(string(r))
	return nil
}

//TODO Stars a repository.
func (d *DockerApi) StarRepository(username, name string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/stars/", username, name))
	r, err := requests.Post(d.client, pth, map[string]string{}, d.token)
	if err != nil {
		return nil
	}
	log.Print(string(r))
	return nil
}

//TODO
func (d *DockerApi) UnstarRepository(username, name string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/stars", username, name))
	r, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		return nil
	}
	log.Print(string(r))
	return nil
}

//GetUser gets info about the given user.
func (d *DockerApi) GetUser(username string) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("no user given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("users/%s", username))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	var user User
	_ = json.Unmarshal(r, &user)
	return &user, nil
}

//AddCollaborator adds a collaborator to an image
func (d *DockerApi) AddCollaborator(username, name, collaborator string) error {
	username = strings.ToLower(username)
	collaborator = strings.ToLower(collaborator)
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if collaborator == "" {
		return fmt.Errorf("no collaborator given")
	}
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/collaborators", username, name))
	_, err := requests.Post(d.client, pth, map[string]string{
		"user": collaborator,
	}, d.token)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerApi) GetUsername() string {
	return d.username
}

func (d *DockerApi) GetToken() string {
	return d.token
}
