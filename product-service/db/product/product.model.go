package product

import "time"

type Product struct {
	ID           string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Name         string `gorm:"index" json:"name"`
	Description  string `json:"description"`
	ImageUrl     string `json:"imageUrl"`
	SearchVector string `gorm:"type:tsvector"`
	SKU          string `gorm:"unique" json:"SKU"`
	Price        int    `json:"price"`
	Currency     string `json:"currency"`

	CategoryId string `json:"categoryId"`
	AdminId    string `json:"adminId"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
