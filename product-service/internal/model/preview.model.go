package model

import (
	"time"
)

type Preview struct {
	ID        string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Position  int    `json:"position"`
	Extension string `json:"extension"`
	Alt       string `json:"alt"`
	Mime      string `json:"mime"`
	Path      string `json:"-"`

	ProductId string    `json:"productId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
