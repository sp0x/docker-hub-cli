package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func authenticateRequest(req *http.Request, token string) {
	req.Header.Add("Authorization", "JWT "+token)
}

func setupHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "DockerHubCli 0.1")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/json")
	//If we request gzip, we have to manually gunzip it.
	//req.Header.Add("Accept-Encoding", "gzip")
}

func Post(client *http.Client, route string, objData interface{}, token string) ([]byte, error) {
	if client == nil {
		return []byte{}, errors.New("null transport client")
	}
	data, err := json.Marshal(objData)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", route, buff)
	setupHeaders(req)
	if token != "" {
		authenticateRequest(req, token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return body, errors.New(strconv.Itoa(res.StatusCode))
	}
	return body, err
}

func Put(client *http.Client, route string, objData interface{}, token string) ([]byte, error) {
	if client == nil {
		return []byte{}, errors.New("null transport client")
	}
	data, err := json.Marshal(objData)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(data)
	req, _ := http.NewRequest("PUT", route, buff)
	setupHeaders(req)
	if token != "" {
		authenticateRequest(req, token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return body, errors.New(strconv.Itoa(res.StatusCode))
	}
	return body, err
}

func Patch(client *http.Client, route string, objData interface{}, token string) ([]byte, error) {
	if client == nil {
		return []byte{}, errors.New("null transport client")
	}
	data, err := json.Marshal(objData)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(data)
	req, _ := http.NewRequest("PATCH", route, buff)
	setupHeaders(req)
	if token != "" {
		authenticateRequest(req, token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return body, errors.New(strconv.Itoa(res.StatusCode))
	}
	return body, err
}

func Get(client *http.Client, route string, token string) ([]byte, error) {
	if client == nil {
		return []byte{}, errors.New("null transport client")
	}
	req, _ := http.NewRequest("GET", route, nil)
	setupHeaders(req)
	if token != "" {
		authenticateRequest(req, token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return body, errors.New(strconv.Itoa(res.StatusCode))
	}
	return body, err
}

func GetWithHeaders(client *http.Client, url string, headers map[string]string) ([]byte, error) {
	if client == nil {
		return []byte{}, errors.New("null transport client")
	}
	req, _ := http.NewRequest("GET", url, nil)
	//setupHeaders(req)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return body, errors.New(strconv.Itoa(res.StatusCode))
	}
	return body, err
}

func Delete(client *http.Client, route string, token string) ([]byte, error) {
	if client == nil {
		return []byte{}, errors.New("null transport client")
	}
	req, _ := http.NewRequest("DELETE", route, nil)
	setupHeaders(req)
	if token != "" {
		authenticateRequest(req, token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return body, errors.New(strconv.Itoa(res.StatusCode))
	}
	return body, err
}
