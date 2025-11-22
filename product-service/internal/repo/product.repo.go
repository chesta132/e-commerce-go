package repo

import (
	"context"
	"product-service/internal/model"
	"strings"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"github.com/chesta132/e-commerce-go/shared/sslicelib"
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

func (r *Product) DB() *gorm.DB {
	return r.db
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

func (r *Product) SearchByKeyword(ctx context.Context, keyword string, opt SearchByKeywordOptions) (products []model.Product, notFoundCatIds []string, err error) {
	subQuery := r.SearchQuery(keyword).WithContext(ctx)

	query := r.db.
		Preload("Meta").
		Table("(?) as sub", subQuery).
		Where("rank > 0.2")

	if len(opt.CategoryIds) > 0 {
		var existingIds []string
		err := r.db.Model(&model.Category{}).
			Where("id IN ?", opt.CategoryIds).
			Pluck("id", &existingIds).Error

		if err != nil {
			return nil, nil, err
		}

		if len(existingIds) < len(opt.CategoryIds) {
			notFoundCatIds = sslicelib.FilterNotInSlice(opt.CategoryIds, existingIds)
		}

		if len(existingIds) == 0 {
			return nil, notFoundCatIds, nil
		}

		query = query.
			Joins("JOIN product_categories pc ON pc.product_id = sub.id").
			Joins("JOIN categories c ON c.id = pc.category_id").
			Where("c.id IN ?", existingIds)
	}

	err = query.
		Order("rank DESC").
		Offset(opt.Offset).
		Limit(opt.Limit).
		Find(&products).Error

	return products, notFoundCatIds, err
}

func (r *Product) CreateOne(ctx context.Context, p *model.Product) error {
	return gorm.G[model.Product](r.db).Create(ctx, p)
}

func (r *Product) FindFirst(ctx context.Context, where []squery.Where) (model.Product, error) {
	q, v := squery.BuildWhere(where)
	return gorm.G[model.Product](r.db).Where(q, v...).First(ctx)
}

func (r *Product) UpdateOne(ctx context.Context, where []squery.Where, update model.Product) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Product](r.db).Where(q, v...).Updates(ctx, update)
	return err
}

func (r *Product) DeleteOne(ctx context.Context, where []squery.Where) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Product](r.db).Where(q, v...).Delete(ctx)
	return err
}

func (r *Product) DeleteCategories(ctx context.Context, id string, catIds []string) error {
	product := model.Product{ID: id}
	var categories []model.Category
	for _, id := range catIds {
		categories = append(categories, model.Category{ID: id})
	}

	return r.db.
		Model(&product).
		Association("Categories").
		Delete(&categories)
}
