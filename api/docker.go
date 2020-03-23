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

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Gravatar string `json:"gravatar_url"`
}

func (d *DockerApi) getRoute(p string) string {
	return joinURL(d.routeBase, p)
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

// login logs in the user and remembers a token to use for authenticated commands.
func (d *DockerApi) Login(username, password string) error {
	loginPath := d.getRoute("users/login")
	r, err := requests.Post(d.client, loginPath, map[string]string{"username": username, "password": password}, d.token)
	if err != nil {
		log.Error(err)
		return err
	}
	var rawmap map[string]json.RawMessage
	err = json.Unmarshal(r, &rawmap)
	if err != nil {
		return err
	}
	tokenStr := ""
	_ = json.Unmarshal(rawmap["token"], &tokenStr)
	d.token = tokenStr
	return nil
}

func (d *DockerApi) Logout() error {
	logoutPath := d.getRoute("logout")
	_, err := requests.Post(d.client, logoutPath, nil, d.token)
	return err
}

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

func (d *DockerApi) GetMyUser() (*User, error) {
	pth := d.getRoute("user")
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal(r, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
