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

func (r *Category) FindManyByIds(ctx context.Context, ids []string) ([]model.Category, error) {
	return gorm.G[model.Category](r.db).Where("id IN ?", ids).Find(ctx)
}

func (r *Category) FindById(ctx context.Context, id string) (model.Category, error) {
	return gorm.G[model.Category](r.db).Where("id = ?", id).First(ctx)
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
