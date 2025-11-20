package service

import (
	"context"
	"os"
	"product-service/config"
	"product-service/internal/lib/thumbnaillib"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"github.com/labstack/echo/v4"
)

type Thumbnail struct {
	tr *repo.Thumbnail
}

type EchoThumbnail struct {
	Thumbnail
	ctx context.Context
	c   echo.Context
}

func NewThumbnail(repoThumbnail *repo.Thumbnail) *Thumbnail {
	return &Thumbnail{repoThumbnail}
}

func (s *Thumbnail) AttachEcho(c echo.Context) *EchoThumbnail {
	return &EchoThumbnail{Thumbnail: *s, ctx: c.Request().Context(), c: c}
}

func (s *EchoThumbnail) CreateThumbnail(thumbnail *model.Thumbnail, content []byte) error {
	if err := s.tr.CreateOne(s.ctx, thumbnail); err != nil {
		return err
	}
	return s.tr.WriteFile(thumbnail.Path, content)
}

func (s *EchoThumbnail) GetThumbnail(id, prodId string) (model.Thumbnail, []byte, error) {
	thumbnail, err := s.tr.FindFirst(s.ctx, []squery.Where{{Name: "id", Value: id}, {Name: "product_id", Value: prodId}})
	if err != nil {
		return model.Thumbnail{}, nil, err
	}

	content, err := s.tr.ReadFile(thumbnail.Path)
	if err != nil {
		return model.Thumbnail{}, nil, err
	}

	return thumbnail, content, nil
}

func (s *EchoThumbnail) UpdateThumbnail(thumbnail *model.Thumbnail, content []byte) error {
	if err := s.tr.WriteFile(thumbnail.Path, content); err != nil {
		return err
	}
	return s.tr.UpdateOne(s.ctx, []squery.Where{{Name: "id", Value: thumbnail.ID}}, *thumbnail)
}

func (s *EchoThumbnail) DeleteThumbnail(id, prodId string) error {
	thumbnail, err := s.tr.FindFirst(s.ctx, []squery.Where{{Name: "id", Value: id}, {Name: "product_id", Value: prodId}})
	if err != nil {
		return err
	}

	if err := s.tr.DeleteFile(thumbnail.Path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return s.tr.DeleteOne(s.ctx, []squery.Where{{Name: "id", Value: thumbnail.ID}})
}

func (s *Thumbnail) ReadDefaultThumbnail() (model.Thumbnail, []byte, error) {
	content, err := s.tr.ReadFile(config.DEFAULT_THUMBNAIL_PATH)
	return thumbnaillib.GetDefaultThumbnail(), content, err
}
