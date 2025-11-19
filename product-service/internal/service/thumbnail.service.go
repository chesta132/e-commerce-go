package service

import (
	"os"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/thumbnaillib"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared/sslicelib"
	"github.com/google/uuid"
)

type Thumbnail struct {
	tr *repo.Thumbnail
}

func NewThumbnail(repoThumbnail *repo.Thumbnail) *Thumbnail {
	return &Thumbnail{repoThumbnail}
}

func (s *Thumbnail) GetMeta(projectId string) (model.ThumbnailMeta, error) {
	metaPath := thumbnaillib.GetMetaPath(projectId)
	return s.tr.ReadMeta(metaPath)
}

func (s *Thumbnail) CreateThumbnail(thumbnail *model.Thumbnail, projectId string, content []byte) error {
	meta, err := s.GetMeta(projectId)
	if err != nil {
		return err
	}

	thumbnail.ID = uuid.NewString()
	thumbnail.Position = len(meta.Thumbnails) + 1
	thumbnail.Path = thumbnaillib.GetFilePath(thumbnail, projectId)

	if err := s.tr.WriteFile(thumbnail.Path, content); err != nil {
		return err
	}

	metaPath := thumbnaillib.GetMetaPath(projectId)
	if err := s.tr.WriteMeta(metaPath, &meta); err != nil {
		s.tr.DeleteFile(thumbnail.Path)
		return err
	}

	meta.Thumbnails = append(meta.Thumbnails, thumbnail)
	thumbnaillib.Reindex(&meta)

	return nil
}

func (s *Thumbnail) GetThumbnail(id string, projectId string) (*model.Thumbnail, []byte, error) {
	meta, err := s.GetMeta(projectId)
	if err != nil {
		return nil, nil, err
	}

	thumbnail, ok := s.tr.FindThumbnailById(meta.Thumbnails, id)
	if !ok {
		return nil, nil, errorlib.ErrThumbnailNotFound
	}

	content, err := s.tr.ReadFile(thumbnail.Path)
	if err != nil {
		return nil, nil, err
	}

	return thumbnail, content, nil
}

func (s *Thumbnail) UpdateThumbnail(thumbnail *model.Thumbnail, projectId string, content []byte) error {
	meta, err := s.GetMeta(projectId)
	if err != nil {
		return err
	}

	if err := s.tr.WriteFile(thumbnail.Path, content); err != nil {
		return err
	}

	meta.Thumbnails = sslicelib.Map(meta.Thumbnails, func(index int, item *model.Thumbnail) *model.Thumbnail {
		if item.ID == thumbnail.ID {
			return thumbnail
		}
		return item
	})

	return nil
}

func (s *Thumbnail) DeleteThumbnail(id, projectId string) error {
	meta, err := s.GetMeta(projectId)
	if err != nil {
		return err
	}

	thumbnail, ok := s.tr.FindThumbnailById(meta.Thumbnails, id)
	if !ok {
		return errorlib.ErrThumbnailNotFound
	}

	if err := s.tr.DeleteFile(thumbnail.Path); err != nil && !os.IsNotExist(err) {
		return err
	}

	newThumbnails := make([]*model.Thumbnail, 0)
	for _, t := range meta.Thumbnails {
		if t.ID != id {
			newThumbnails = append(newThumbnails, t)
		}
	}
	meta.Thumbnails = newThumbnails
	thumbnaillib.Reindex(&meta)

	metaPath := thumbnaillib.GetMetaPath(projectId)
	return s.tr.WriteMeta(metaPath, &meta)
}
