package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func NewApi() *DockerApi {
	d := &DockerApi{}
	version := "2"
	transport := &http.Transport{
		DisableCompression: false,
	}
	d.client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	d.apiVersion = version
	d.routeBase = fmt.Sprintf("https://hub.docker.com/v%s", version)
	return d
}

type DockerApi struct {
	client     *http.Client
	apiVersion string
	routeBase  string
	token      string
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func (d *DockerApi) getRoute(p string) string {
	return joinURL(d.routeBase, p)
}

// login logs in the user and remembers a token to use for authenticated commands.
func (d *DockerApi) login(username, password string) error {
	loginPath := d.getRoute("users/login")
	r, err := post(d, loginPath, map[string]string{"username": username, "password": password}, false)
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
	json.Unmarshal(rawmap["token"], &tokenStr)
	d.token = tokenStr
	return nil
}

func (d *DockerApi) logout() error {
	logoutPath := d.getRoute("logout")
	_, err := post(d, logoutPath, nil, true)
	return err
}

func (d *DockerApi) getBuildSettings(username string, name string) string {
	username = strings.ToLower(username)
	path := d.getRoute(fmt.Sprintf("repositories/%s/%s/autobuild", username, name))
	r, err := get(d, path, false)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(r)
}

func authenticateRequest(d *DockerApi, req *http.Request) {
	req.Header.Add("authorization", "JWT "+d.token)
}

func jsonifyRequest(req *http.Request) {
	req.Header.Add("content-type", "application/json")
	//If we request gzip, we have to manually gunzip it.
	//req.Header.Add("Accept-Encoding", "gzip")
}

func post(d *DockerApi, route string, objData interface{}, withAuth bool) ([]byte, error) {
	if d.client == nil {
		return []byte{}, errors.New("null transport client")
	}
	data, err := json.Marshal(objData)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", route, buff)
	req.Header.Add("cache-control", "no-cache")
	jsonifyRequest(req)
	if withAuth {
		if d.token == "" {
			return nil, fmt.Errorf("not authenticated")
		}
		authenticateRequest(d, req)
	}
	res, err := d.client.Do(req)
	if err != nil {
		log.Printf("order error: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func get(d *DockerApi, route string, withAuth bool) ([]byte, error) {
	if d.client == nil {
		return []byte{}, errors.New("null transport client")
	}
	req, _ := http.NewRequest("GET", route, nil)
	req.Header.Add("cache-control", "no-cache")
	jsonifyRequest(req)
	if withAuth {
		if d.token == "" {
			return nil, fmt.Errorf("not authenticated")
		}
		authenticateRequest(d, req)
	}

	res, err := d.client.Do(req)
	if err != nil {
		log.Printf("order error: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}
