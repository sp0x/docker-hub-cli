package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/docker-hub-cli/requests"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Gravatar string `json:"gravatar_url"`
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
	d.username = username
	return nil
}

// Logout  of the current user
func (d *DockerApi) Logout() error {
	logoutPath := d.getRoute("logout")
	_, err := requests.Post(d.client, logoutPath, nil, d.token)
	return err
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
