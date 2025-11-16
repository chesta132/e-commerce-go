package model

import "time"

type StockMovement struct {
	ID          string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	ProductId   string `json:"productId"`
	Quantity    int    `json:"quantity"`
	Information string `json:"Information"`
	Type        string `json:"type"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
