package repo

import (
	"context"
	"product-service/internal/model"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"gorm.io/gorm"
)

type Category struct {
	db *gorm.DB
}

func NewCategory(db *gorm.DB) *Category {
	return &Category{db}
}

func (r *Category) DB() *gorm.DB {
	return r.db
}

func (r *Category) CreateOne(ctx context.Context, category *model.Category) error {
	return gorm.G[model.Category](r.db).Create(ctx, category)
}

func (r *Category) CreateMany(ctx context.Context, categories *[]model.Category) error {
	return r.db.WithContext(ctx).Create(categories).Error
}

func (r *Category) FindMany(ctx context.Context, where []squery.Where) ([]model.Category, error) {
	q, v := squery.BuildWhere(where)
	return gorm.G[model.Category](r.db).Where(q, v...).Find(ctx)
}

func (r *Category) FindFirst(ctx context.Context, where []squery.Where) (model.Category, error) {
	q, v := squery.BuildWhere(where)
	return gorm.G[model.Category](r.db).Where(q, v...).First(ctx)
}

func (r *Category) UpdateOne(ctx context.Context, where []squery.Where, update model.Category) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Category](r.db).Where(q, v...).Updates(ctx, update)
	return err
}

func (r *Category) DeleteOne(ctx context.Context, where []squery.Where) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Category](r.db).Where(q, v...).Delete(ctx)
	return err
}

func (s *Category) FindProductsByCategoryIds(ids []string) ([]model.Product, error) {
	var products []model.Product
	err := s.db.Joins("JOIN product_categories pc ON pc.product_id = products.id").
		Where("pc.category_id IN ?", ids).
		Preload("Categories").
		Find(&products).Error
	return products, err
}
