package service

import (
	"context"
	_ "image/png"
	"product-service/internal/model"
	"product-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Image struct {
	ir *repo.Image
}

type EchoImage struct {
	Image
	ctx context.Context
	c   echo.Context
}

func NewImage(repo *repo.Image) *Image {
	return &Image{repo}
}

func (s *Image) AttachEcho(c echo.Context) *EchoImage {
	return &EchoImage{Image: *s, ctx: c.Request().Context(), c: c}
}

func (s *EchoImage) CreateOne(img *model.Image) error {
	return s.ir.CreateOne(s.ctx, img)
}
