package repo

import (
	"context"
	"product-service/internal/model"

	"gorm.io/gorm"
)

type Image struct {
	db *gorm.DB
}

func NewImage(db *gorm.DB) *Image {
	return &Image{db}
}

func (r *Image) CreateOne(ctx context.Context, img *model.Image) error {
	return gorm.G[model.Image](r.db).Create(ctx, img)
}
