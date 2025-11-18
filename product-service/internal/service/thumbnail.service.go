package service

import (
	"context"
	"product-service/internal/repo"

	"github.com/labstack/echo/v4"
)

type Thumbnail struct {
	ir *repo.Thumbnail
}

type EchoThumbnail struct {
	Thumbnail
	ctx context.Context
	c   echo.Context
}

func NewThumbnail(repo *repo.Thumbnail) *Thumbnail {
	return &Thumbnail{repo}
}

func (s *Thumbnail) AttachEcho(c echo.Context) *EchoThumbnail {
	return &EchoThumbnail{Thumbnail: *s, ctx: c.Request().Context(), c: c}
}

func FindOne(thumbId, prodId string) {
	
}