package repo

import (
	"context"
	"product-service/internal/model"
	"strings"

	"gorm.io/gorm"
)

type Product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{db}
}

func (r *Product) SearchByKeyword(ctx context.Context, keyword string, offset, limit int) ([]model.Product, error) {
	var products []model.Product

	prefixQuery := strings.ReplaceAll(keyword, " ", ":* & ") + ":*"

	subQuery := r.db.WithContext(ctx).
		Table("products").
		Select(`products.*, 
			COALESCE(ts_rank(search_vector, plainto_tsquery('english', ?)), 0) * 10 +
			COALESCE(similarity(name, ?), 0) * 5 +
			COALESCE(similarity(description, ?), 0) * 2 as rank`,
			keyword, keyword, keyword).
		Where(`
			search_vector @@ to_tsquery('english', ?) OR
			search_vector @@ plainto_tsquery('english', ?) OR
			name ILIKE ? OR
			description ILIKE ? OR
			similarity(name, ?) > 0.2 OR
			similarity(description, ?) > 0.2`,
			prefixQuery, keyword, "%"+keyword+"%", "%"+keyword+"%", keyword, keyword)

	err := r.db.WithContext(ctx).
		Preload("Meta").
		Table("(?) as sub", subQuery).
		Where("rank > 0.2").
		Order("rank DESC").
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	return products, err
}

func (r *Product) CreateOne(ctx context.Context, p *model.Product) error {
	return gorm.G[model.Product](r.db).Create(ctx, p)
}

func (r *Product) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
