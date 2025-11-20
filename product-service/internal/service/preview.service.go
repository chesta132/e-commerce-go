package service

import (
	"context"
	"os"
	"product-service/config"
	"product-service/internal/lib/previewlib"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"github.com/labstack/echo/v4"
)

type Preview struct {
	tr *repo.Preview
}

type EchoPreview struct {
	Preview
	ctx context.Context
	c   echo.Context
}

func NewPreview(repoPreview *repo.Preview) *Preview {
	return &Preview{repoPreview}
}

func (s *Preview) AttachEcho(c echo.Context) *EchoPreview {
	return &EchoPreview{Preview: *s, ctx: c.Request().Context(), c: c}
}

func (s *EchoPreview) CreatePreview(preview *model.Preview, content []byte) error {
	if err := s.tr.CreateOne(s.ctx, preview); err != nil {
		return err
	}
	return s.tr.WriteFile(preview.Path, content)
}

func (s *EchoPreview) GetPreview(id, prodId string) (model.Preview, []byte, error) {
	preview, err := s.tr.FindFirst(s.ctx, []squery.Where{{Name: "id", Value: id}, {Name: "product_id", Value: prodId}})
	if err != nil {
		return model.Preview{}, nil, err
	}

	content, err := s.tr.ReadFile(preview.Path)
	if err != nil {
		return model.Preview{}, nil, err
	}

	return preview, content, nil
}

func (s *EchoPreview) UpdatePreview(preview *model.Preview, content []byte) error {
	if err := s.tr.WriteFile(preview.Path, content); err != nil {
		return err
	}
	return s.tr.UpdateOne(s.ctx, []squery.Where{{Name: "id", Value: preview.ID}}, *preview)
}

func (s *EchoPreview) DeletePreview(id, prodId string) error {
	preview, err := s.tr.FindFirst(s.ctx, []squery.Where{{Name: "id", Value: id}, {Name: "product_id", Value: prodId}})
	if err != nil {
		return err
	}

	if err := s.tr.DeleteFile(preview.Path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return s.tr.DeleteOne(s.ctx, []squery.Where{{Name: "id", Value: preview.ID}})
}

func (s *Preview) ReadDefaultPreview() (model.Preview, []byte, error) {
	content, err := s.tr.ReadFile(config.DEFAULT_PREVIEW_PATH)
	return previewlib.GetDefaultPreview(), content, err
}
