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

type Thumbnail struct {
	db *gorm.DB
}

func NewThumbnail(db *gorm.DB) *Thumbnail {
	return &Thumbnail{db}
}

func (r *Thumbnail) WriteFile(path string, content []byte) error {
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

func (r *Thumbnail) ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (r *Thumbnail) DeleteFile(path string) error {
	return os.Remove(path)
}

func (r *Thumbnail) CreateOne(ctx context.Context, thumbnail *model.Thumbnail) error {
	return gorm.G[model.Thumbnail](r.db).Create(ctx, thumbnail)
}

func (r *Thumbnail) CreateMany(ctx context.Context, thumbnail *[]model.Thumbnail) error {
	return gorm.G[model.Thumbnail](r.db).CreateInBatches(ctx, thumbnail, len(*thumbnail))
}

func (r *Thumbnail) FindFirst(ctx context.Context, where []squery.Where) (model.Thumbnail, error) {
	q, v := squery.BuildWhere(where)
	return gorm.G[model.Thumbnail](r.db).Where(q, v...).First(ctx)
}

func (r *Thumbnail) UpdateOne(ctx context.Context, where []squery.Where, thumbnail model.Thumbnail) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Thumbnail](r.db).Where(q, v...).Updates(ctx, thumbnail)
	return err
}

func (r *Thumbnail) DeleteOne(ctx context.Context, where []squery.Where) error {
	q, v := squery.BuildWhere(where)
	_, err := gorm.G[model.Thumbnail](r.db).Where(q, v...).Delete(ctx)
	return err
}
