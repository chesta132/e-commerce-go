package repo

import (
	"context"
	"product-service/internal/model"

	"gorm.io/gorm"
)

type Product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{db}
}

func (r *Product) SearchByKeyword(ctx context.Context, keyword string, offset, limit int) ([]model.Product, error) {
	return gorm.G[model.Product](r.db).
		Select("*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", keyword). // as rank for order
		Where("search_vector @@ plainto_tsquery('english', ?)", keyword).
		Order("rank DESC").
		Offset(offset).
		Limit(limit).
		Find(ctx)
}

func (r *Product) CreateOne(ctx context.Context, p *model.Product) error {
	return gorm.G[model.Product](r.db).Create(ctx, p)
}

func (r *Product) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
