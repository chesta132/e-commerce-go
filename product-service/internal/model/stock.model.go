package model

import "time"

type Stock struct {
	ID        string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Reserved  int    `json:"reserved"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
