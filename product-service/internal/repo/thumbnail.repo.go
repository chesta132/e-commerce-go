package repo

import (
	"encoding/json"
	"io"
	"os"
	"product-service/internal/lib/errorlib"
	"product-service/internal/model"

	"github.com/chesta132/e-commerce-go/shared/sslicelib"
)

type Thumbnail struct{}

func NewThumbnail() *Thumbnail {
	return &Thumbnail{}
}

func (r *Thumbnail) WriteMeta(path string, meta *model.ThumbnailMeta) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	smeta, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(smeta)
	return err
}

func (r *Thumbnail) ReadMeta(path string) (model.ThumbnailMeta, error) {
	var meta model.ThumbnailMeta
	metaByte, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return model.ThumbnailMeta{}, errorlib.ErrThumbnailMetaNotFound
	}
	if err != nil {
		return model.ThumbnailMeta{}, err
	}
	err = json.Unmarshal(metaByte, &meta)
	return meta, err
}

func (r *Thumbnail) WriteFile(path string, content []byte) error {
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

func (r *Thumbnail) FindThumbnailById(thumbnails []*model.Thumbnail, id string) (*model.Thumbnail, bool) {
	return sslicelib.Find(thumbnails, func(index int, item *model.Thumbnail) bool {
		return item.ID == id
	})
}

func (r *Thumbnail) DeleteFile(path string) error {
	return os.Remove(path)
}
