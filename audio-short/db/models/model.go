package models

import "encoding/json"

type Creator struct {
	Name  string `string:"name,omitempty"`
	Email string `string:"email,omitempty"`
}

// short schema of the short table
type Short struct {
	ID          int `string:"Id,omitempty"`
	Title       string `string:"Title,omitempty"`
	Description string `string:"Description,omitempty"`
	Category    string `string:"Category,omitempty"`
	FileUrl     string `string:"FileUrl,omitempty"`
	Creator     json.RawMessage   `json:"Creator,omitempty"`
}