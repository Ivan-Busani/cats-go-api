package model

import (
	"encoding/json"
	"time"
)

type Cat struct {
	ID        int             `json:"id"`
	CatID     string          `json:"cat_id"`
	URL       string          `json:"url"`
	Width     int             `json:"width"`
	Height    int             `json:"height"`
	Breeds    json.RawMessage `json:"breeds"`
	APIUsed   string          `json:"api_used"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

type SaveCatInput struct {
	CatID  string          `json:"cat_id"`
	URL    string          `json:"url"`
	Width  int             `json:"width"`
	Height int             `json:"height"`
	Breeds json.RawMessage `json:"breeds"`
}
