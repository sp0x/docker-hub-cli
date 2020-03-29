package api

import (
	"encoding/json"
	"fmt"
	"github.com/sp0x/docker-hub-cli/requests"
	"strings"
	"time"
)

type Webhook struct {
	Id                  *int         `json:"id"`
	Name                string       `json:"name"`
	Active              *bool        `json:"active"`
	ExpectFinalCallback *bool        `json:"expect_final_callback"`
	Creator             *string      `json:"creator"`
	Hooks               []WebhookUrl `json:"hooks"`
	HookUrl             *string      `json:"hook_url"`
	Created             *time.Time   `json:"created"`
	LastUpdated         *time.Time   `json:"last_updated"`
	LastUpdater         *string      `json:"last_updater"`
}

type WebhookUrl struct {
	Name string `json:"name"`
	Url  string `json:"hook_url"`
}

func (wh *Webhook) GetWebhookUrl() *string {
	if wh.Hooks != nil && len(wh.Hooks) > 0 {
		return &wh.Hooks[0].Url
	}
	return wh.HookUrl
}

//DeleteAllWebhooks deletes all webhooks for a given repository
func (d *DockerApi) DeleteAllWebhooks(username, name string) error {
	hooks, err := d.GetWebhooks(username, name, 0, 0)
	if err != nil {
		return err
	}
	for _, h := range hooks {
		err := d.DeleteWebhook(username, name, h.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteWebhook deletes a webhook
func (d *DockerApi) DeleteWebhook(username, name, webhookName string) error {
	if username == "" {
		return fmt.Errorf("no user given")
	}
	if name == "" {
		return fmt.Errorf("no image name given")
	}
	if webhookName == "" {
		return fmt.Errorf("no webhookName given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/webhook_pipeline/%s/", username, name, webhookName)) + "/"
	r, err := requests.Delete(d.client, pth, d.token)
	if err != nil {
		if r != nil {
			return fmt.Errorf(parseError(r))
		} else {
			return err
		}
	}
	return nil
}

//GetWebhooks Gets the webhooks for a repository you own.
func (d *DockerApi) GetWebhooks(username, name string, pageSize, page int) ([]Webhook, error) {
	if username == "" || username == "_" {
		username = "library"
	}
	if name == "" {
		return nil, fmt.Errorf("no image name given")
	}
	username = strings.ToLower(username)
	name = strings.ToLower(name)
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/webhook_pipeline?page_size=%v&page=%v", username, name, pageSize, page))
	r, err := requests.Get(d.client, pth, d.token)
	if err != nil {
		return nil, err
	}
	var searchRes SearchResult
	err = json.Unmarshal(r, &searchRes)
	if err != nil {
		return nil, err
	}
	var webhooks []Webhook
	err = json.Unmarshal(searchRes.Results, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

//CreateWebhook Creates a webhook for the given username and repository.
func (d *DockerApi) CreateWebhook(username, name, webhookName string, url string) (*Webhook, error) {
	if username == "" {
		return nil, fmt.Errorf("no user given")
	}
	if name == "" {
		return nil, fmt.Errorf("no image name given")
	}
	if webhookName == "" {
		return nil, fmt.Errorf("no webhookName given")
	}
	username = strings.ToLower(username)
	pth := d.getRoute(fmt.Sprintf("repositories/%s/%s/webhook_pipeline/", username, name)) + "/"
	data := map[string]interface{}{
		"name":                  webhookName,
		"expect_final_callback": false,
		"webhooks": []map[string]string{
			{
				"name": webhookName, "hook_url": url,
			}},
		"registry": "registry-1.docker.io",
	}
	r, err := requests.Post(d.client, pth, data, d.token)
	if err != nil {
		if r != nil {
			return nil, fmt.Errorf(parseError(r))
		} else {
			return nil, err
		}
	}
	var hook Webhook
	err = json.Unmarshal(r, &hook)
	if err != nil {
		return nil, err
	}
	return &hook, nil
}
