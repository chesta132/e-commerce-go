package repo

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"product-service/internal/model"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"gorm.io/gorm"
)

type Preview struct {
	db *gorm.DB
}

func NewPreview(db *gorm.DB) *Preview {
	return &Preview{db}
}

func (r *Preview) WriteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

func (r *Preview) ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (r *Preview) DeleteFile(path string) error {
	return os.Remove(path)
}

func (r *Preview) CreateOne(ctx context.Context, preview *model.Preview) error {
	return gorm.G[model.Preview](r.db).Create(ctx, preview)
}

func (r *Preview) CreateMany(ctx context.Context, preview *[]model.Preview) error {
	return gorm.G[model.Preview](r.db).CreateInBatches(ctx, preview, len(*preview))
}

func (r *Preview) FindFirst(ctx context.Context, where []squery.Where) (model.Preview, error) {
	q, v := squery.BuildWhere(where)
	return gorm.G[model.Preview](r.db).Where(q, v...).First(ctx)
}

func (r *Preview) UpdateOne(ctx context.Context, where []squery.Where, preview model.Preview) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Preview](r.db).Where(q, v...).Updates(ctx, preview)
	return err
}

func (r *Preview) DeleteOne(ctx context.Context, where []squery.Where) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Preview](r.db).Where(q, v...).Delete(ctx)
	return err
}
