package model

import (
	"time"
)

type Product struct {
	ID           string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Name         string `gorm:"index" json:"name"`
	Description  string `gorm:"index" json:"description"`
	SearchVector string `gorm:"type:tsvector;->" json:"-"`
	SKU          string `gorm:"unique"`

	MetaId string      `gorm:"unique" json:"metaId"`
	Meta   ProductMeta `gorm:"foreignKey:MetaId;constraint:OnDelete:CASCADE" json:"meta,omitzero"`

	Previews   []Preview  `gorm:"constraint:OnDelete:CASCADE" json:"previews,omitempty"`
	Categories []Category `gorm:"many2many:product_categories" json:"categories,omitempty"`
	AdminId    string     `json:"adminId" gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	PreviewIds  []string `json:"previewIds,omitempty" gorm:"-"`
	CategoryIds []string `json:"categoryIds,omitempty" gorm:"-"`
}

type ProductMeta struct {
	ID            string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Price         int    `json:"price"`
	Currency      string `json:"currency"`
	TotalQuantity int    `json:"quantity"`
	TotalReserved int    `json:"reserved"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type CreateProductPayload struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Price       int      `json:"price" validate:"required"`
	Currency    string   `json:"currency" validate:"required"`
	CategoryIds []string `json:"categoryIds" validate:"required"`
}

type DeleteCategoryRelationPayload struct {
	CatIds []string `json:"categoryIds" validate:"required"`
}
