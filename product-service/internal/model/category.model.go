package model

import "time"

type Category struct {
	ID          string    `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"many2many:product_categories;" json:"products"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
