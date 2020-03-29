package api

import "encoding/json"

type SearchResult struct {
	Count int `json:"count"`
	//An url to the next search page
	Next *string `json:"next"`
	//An url to the previous search page
	Previous *string `json:"previous"`
	Results  json.RawMessage
}
