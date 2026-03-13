package model

import (
	"encoding/json"
	"time"
)

type Cat struct {
	ID        int             `json:"id"         db:"id"`
	CatID     string          `json:"cat_id"     db:"cat_id"`
	URL       string          `json:"url"        db:"url"`
	Width     int             `json:"width"      db:"width"`
	Height    int             `json:"height"     db:"height"`
	Breeds    json.RawMessage `json:"breeds"     db:"breeds"`
	APIUsed   *string         `json:"api_used"   db:"api_used"`
	CreatedAt *time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at" db:"updated_at"`
}

type SaveCatInput struct {
	CatID  string          `json:"cat_id"`
	URL    string          `json:"url"`
	Width  int             `json:"width"`
	Height int             `json:"height"`
	Breeds json.RawMessage `json:"breeds"`
}
