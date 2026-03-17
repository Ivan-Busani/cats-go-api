package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type Cat struct {
	ID        int            `json:"id"         gorm:"column:id;primaryKey;autoIncrement"`
	CatID     string         `json:"cat_id"     gorm:"column:cat_id"`
	URL       string         `json:"url"        gorm:"column:url"`
	Width     int            `json:"width"      gorm:"column:width"`
	Height    int            `json:"height"     gorm:"column:height"`
	Breeds    datatypes.JSON `json:"breeds"     gorm:"column:breeds"`
	APIUsed   *string        `json:"api_used"   gorm:"column:api_used"`
	CreatedAt *time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"column:updated_at"`
}

func (Cat) TableName() string {
	return "cats"
}

type SaveCatInput struct {
	CatID  string          `json:"cat_id"`
	URL    string          `json:"url"`
	Width  int             `json:"width"`
	Height int             `json:"height"`
	Breeds json.RawMessage `json:"breeds"`
}
