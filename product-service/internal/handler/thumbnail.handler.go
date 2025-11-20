package handler

import (
	"bytes"
	"errors"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/lib/thumbnaillib"
	"product-service/internal/service"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Thumbnail struct {
	svc  *service.Thumbnail
	psvc *service.Product
}

func NewThumbnail(svc *service.Thumbnail, psvc *service.Product) *Thumbnail {
	return &Thumbnail{svc, psvc}
}

func (h *Thumbnail) GetOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	prodId := c.Param("prod-id")
	id := c.Param("id")

	meta, thumb, err := svc.GetThumbnail(id, prodId)
	if err != nil {
		_, t, e := svc.ReadDefaultThumbnail()
		if e != nil {
			return rp.Error(replylib.CodeServerError, e.Error()).Info("Error while read default thumbnail").FailJSON()
		}
		return errorlib.HandleGetThumbnailError(err, t, rp)
	}

	r := bytes.NewReader(thumb)
	return rp.AddHeaders(thumbnaillib.GetHeader(meta)).Success(reply.Stream{Data: r, ContentType: meta.Mime}).OkStream()
}

func (h *Thumbnail) CreateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")
	if _, err := psvc.FindById(prodId); errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeBadRequest, "record: product with requested id not found").FailJSON()
	}

	alt := c.FormValue("alt")
	fh, err := c.FormFile("file")
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}

	if alt == "" {
		alt = thumbnaillib.GetFileName(fh.Filename)
	}

	meta := thumbnaillib.GenerateThumbnail(fh, alt, prodId, true)

	file, fbyte, err := thumbnaillib.ReadByHeader(fh)
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	defer file.Close()
	err = svc.CreateThumbnail(meta, fbyte)
	if err != nil {
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(meta).CreatedJSON()
}

func (h *Thumbnail) UpdateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")
	id := c.Param("id")
	if _, err := psvc.FindById(prodId); errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeBadRequest, "record: product with requested id not found").FailJSON()
	}
	meta, _, err := svc.GetThumbnail(id, prodId)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	alt := c.FormValue("alt")
	meta.Alt = alt
	thumbnaillib.MergeThumbnail(fh, &meta)

	f, b, err := thumbnaillib.ReadByHeader(fh)
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	defer f.Close()

	err = svc.UpdateThumbnail(&meta, b)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}
	return rp.Success(meta).OkJSON()
}
