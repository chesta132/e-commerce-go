package repo

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/thumbnaillib"
	"product-service/internal/model"

	"github.com/chesta132/e-commerce-go/shared/sslicelib"
	"github.com/google/uuid"
)

type Thumbnail struct{}

func NewThumbnail() *Thumbnail {
	return &Thumbnail{}
}

func (r *Thumbnail) StoreMeta(meta *model.ThumbnailMeta) error {
	path := filepath.Join(thumbnaillib.GetDirPath(meta.ProjectId), "meta.json")

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	smeta, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(smeta)
	return err
}

func (r *Thumbnail) ReadMeta(projectId string) (model.ThumbnailMeta, error) {
	var meta model.ThumbnailMeta
	dirPath := thumbnaillib.GetDirPath(projectId)
	metaByte, err := os.ReadFile(filepath.Join(dirPath, "meta.json"))
	if os.IsNotExist(err) {
		return model.ThumbnailMeta{}, errorlib.ErrThumbnailMetaNotFound
	}
	if err != nil {
		return model.ThumbnailMeta{}, err
	}
	err = json.Unmarshal(metaByte, &meta)
	return meta, err
}

func (r *Thumbnail) ReindexAndStore(meta *model.ThumbnailMeta) error {
	thumbnaillib.Reindex(meta)
	return r.StoreMeta(meta)
}

func (r *Thumbnail) CreateThumbnail(thumbnail *model.Thumbnail, projectId string, content []byte) error {
	meta, err := r.ReadMeta(projectId)
	if err != nil {
		return err
	}
	thumbnaillib.Reindex(&meta)
	thumbnail.ID = uuid.NewString()
	thumbnail.Position = len(meta.Thumbnails) + 1
	thumbnail.Path = thumbnaillib.GetFilePath(thumbnail, meta.ProjectId)
	meta.Thumbnails = append(meta.Thumbnails, thumbnail)

	file, err := os.Create(thumbnail.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return r.ReindexAndStore(&meta)
}

func (r *Thumbnail) ReadThumbnail(id string, meta *model.ThumbnailMeta) (*model.Thumbnail, []byte, error) {
	thumbnail, ok := sslicelib.Find(meta.Thumbnails, func(index int, item *model.Thumbnail) bool { return item.ID == id })
	if !ok {
		return nil, nil, errorlib.ErrThumbnailNotFound
	}
	file, err := os.Open(thumbnail.Path)
	if err != nil {
		return nil, nil, err
	}
	fbytes, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	return thumbnail, fbytes, nil
}
