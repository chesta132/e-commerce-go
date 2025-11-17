package model

import "time"

type Category struct {
	ID          string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type CreateCategoryPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}
