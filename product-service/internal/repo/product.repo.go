package repo

import (
	"context"
	"product-service/db/product"

	"gorm.io/gorm"
)

type Product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{db}
}

func (r *Product) SearchByKeyword(ctx context.Context, keyword string, offset, limit int) ([]product.Product, error) {
	return gorm.G[product.Product](r.db).
		Select("*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", keyword). // as rank for order
		Where("search_vector @@ plainto_tsquery('english', ?)", keyword).
		Order("rank DESC").
		Offset(offset).
		Limit(limit).
		Find(ctx)
}
