package api

import (
	"encoding/json"
	"strings"
)

type dockerError struct {
	Error  *string  `json:"error"`
	Name   []string `name:"name"`
	Detail *string  `json:"detail"`
}

func parseError(errb []byte) string {
	var data dockerError
	err := json.Unmarshal(errb, &data)
	if err != nil {
		return "could not parse error"
	}
	if data.Error != nil {
		return *data.Error
	}
	if data.Detail != nil {
		return *data.Detail
	}
	if len(data.Name) > 0 {
		return strings.Join(data.Name, "\n")
	}
	return "unresolved error"
}
