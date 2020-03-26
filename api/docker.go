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
		Jar:       cookies,
	}
	d.cookieJar = cookies

	d.apiVersion = version
	d.routeBase = fmt.Sprintf("https://hub.docker.com/v%s", version)
	return d
}

type DockerApi struct {
	client     *http.Client
	apiVersion string
	routeBase  string
	token      string
	cookieJar  *cookiejar.Jar
}

func (d *DockerApi) getRoute(p string) string {
	return joinURL(d.routeBase, p)
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
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
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s", username, name))
	data := map[string]string{}
	if descLong != "" {
		data["full_description"] = descLong
	}
	if descShort != "" {
		data["description"] = descShort
	}
	r, err := requests.Patch(d.client, pth, data, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(string(r))
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

//GetBuildSettings gets the build settings for an image
func (d *DockerApi) GetBuildSettings(username string, name string) string {
	username = strings.ToLower(username)
	settingsPath := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild", username, name))
	r, err := requests.Get(d.client, settingsPath, d.token)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(r)
}

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

func (d *DockerApi) SaveBuildTag(username, name string, id string, details string) error {

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
func (d *DockerApi) GetTags(username, name string, pageSize, page int) error {
	if username != "" && name == "" {
		name = username
		username = "library "
	}
	if username == "" || username == "_" {
		username = "library"
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/tags?page_size=%v&page=%v", username, name, pageSize, page))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil
	}
	log.Print(string(r))
	return nil
}

type Repository struct {
}

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
	_ = json.Unmarshal(r, &repo)
	return &repo, nil
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

//GetWebhooks Gets the webhooks for a repository you own.
func (d *DockerApi) GetWebhooks(username, name string, pageSize, page int) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/repositories/webhooks?page_size=%v&page=%v", username, name, pageSize, page))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Print(r)
	return nil
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
