package repo

import (
	"context"
	"product-service/internal/model"

	"gorm.io/gorm"
)

type Category struct {
	db *gorm.DB
}

func NewCategory(db *gorm.DB) *Category {
	return &Category{db}
}

func (r *Category) CreateOne(ctx context.Context, category *model.Category) error {
	return gorm.G[model.Category](r.db).Create(ctx, category)
}

func (r *Category) CreateMany(ctx context.Context, categories *[]model.Category) error {
	return r.db.WithContext(ctx).Create(categories).Error
}

func (r *Category) FindManyByIds(ctx context.Context, ids []string) ([]model.Category, error) {
	return gorm.G[model.Category](r.db).Where("id IN ?", ids).Find(ctx)
}
