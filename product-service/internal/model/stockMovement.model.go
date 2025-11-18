package model

import "time"

type StockMovement struct {
	ID          string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Quantity    int    `json:"quantity"`
	Information string `json:"Information"`
	Type        string `json:"type"`

	ProductId string `gorm:"index;constraint:OnDelete:CASCADE" json:"productId"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
