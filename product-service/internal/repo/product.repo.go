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

type SearchByKeywordOptions struct {
	Offset, Limit int
	CategoryIds   []string
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{db}
}

func (r *Product) SearchQuery(keyword string) (query *gorm.DB) {
	pre := strings.ReplaceAll(keyword, " ", ":* & ") + ":*"
	return r.db.
		Table("products").
		Select(`products.*, 
			COALESCE(ts_rank(search_vector, plainto_tsquery('english', ?)), 0) * 10 +
			COALESCE(similarity(products.name, ?), 0) * 5 +
			COALESCE(similarity(products.description, ?), 0) * 2 as rank`,
			keyword, keyword, keyword).
		Where(`
			search_vector @@ to_tsquery('english', ?) OR
			search_vector @@ plainto_tsquery('english', ?) OR
			products.name ILIKE ? OR
			products.description ILIKE ? OR
			similarity(products.name, ?) > 0.2 OR
			similarity(products.description, ?) > 0.2`,
			pre, keyword, "%"+keyword+"%", "%"+keyword+"%", keyword, keyword)
}

func (r *Product) SearchByKeyword(ctx context.Context, keyword string, opt SearchByKeywordOptions) ([]model.Product, error) {
	var products []model.Product

	subQuery := r.SearchQuery(keyword).WithContext(ctx)
	if len(opt.CategoryIds) > 0 {
		subQuery = subQuery.
			Joins("JOIN product_categories pc ON pc.product_id = products.id").
			Joins("JOIN categories c ON c.id = pc.category_id").
			Where("c.id IN ?", opt.CategoryIds)
	}

	err := r.db.
		Preload("Meta").
		Preload("Categories").
		Table("(?) as sub", subQuery).
		Where("rank > 0.2").
		Order("rank DESC").
		Offset(opt.Offset).
		Limit(opt.Limit).
		Find(&products).Error

	return products, err
}

func (r *Product) CreateOne(ctx context.Context, p *model.Product) error {
	return gorm.G[model.Product](r.db).Create(ctx, p)
}

func (r *Product) DB() *gorm.DB {
	return r.db
}
