package handler

import (
	"bytes"
	"errors"
	"product-service/config"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/previewlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/service"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Preview struct {
	svc  *service.Preview
	psvc *service.Product
}

func NewPreview(svc *service.Preview, psvc *service.Product) *Preview {
	return &Preview{svc, psvc}
}

func (h *Preview) GetOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	prodId := c.Param("prod-id")
	id := c.Param("id")

	meta, preview, err := svc.GetPreview(id, prodId)
	if err != nil {
		_, p, e := svc.ReadDefaultPreview()
		if e != nil {
			return rp.Error(replylib.CodeServerError, e.Error()).Info("Error while read default preview").FailJSON()
		}
		return errorlib.HandleGetPreviewError(err, p, rp)
	}

	r := bytes.NewReader(preview)
	return rp.AddHeaders(previewlib.GetHeader(meta)).Success(reply.Stream{Data: r, ContentType: meta.Mime}).OkStream()
}

func (h *Preview) GetInfo(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")

	product, err := psvc.FindByIdWithRelation(prodId, []string{"Previews"})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeBadRequest, "record: product with requested id not found").FailJSON()
	} else if err != nil {
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(product.Previews).OkJSON()
}

func (h *Preview) CreateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")
	product, err := psvc.FindByIdWithRelation(prodId, []string{"Previews"})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeBadRequest, "record: product with requested id not found").FailJSON()
	} else if err != nil {
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}
	if len(product.Previews) >= config.MAX_PREVIEW {
		return rp.Error(replylib.CodeBadRequest, "record: max preview reached").FailJSON()
	}

	alt := c.FormValue("alt")
	fh, err := c.FormFile("file")
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}

	if alt == "" {
		alt = previewlib.GetFileName(fh.Filename)
	}

	meta := previewlib.GeneratePreview(fh, alt, prodId, true)

	file, fbyte, err := previewlib.ReadByHeader(fh)
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	defer file.Close()
	err = svc.CreatePreview(meta, fbyte)
	if err != nil {
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(meta).CreatedJSON()
}

func (h *Preview) UpdateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")
	id := c.Param("id")
	if _, err := psvc.FindById(prodId); errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeBadRequest, "record: product with requested id not found").FailJSON()
	}
	meta, _, err := svc.GetPreview(id, prodId)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	alt := c.FormValue("alt")
	meta.Alt = alt
	previewlib.MergePreview(fh, &meta)

	f, b, err := previewlib.ReadByHeader(fh)
	if err != nil {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	defer f.Close()

	err = svc.UpdatePreview(&meta, b)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}
	return rp.Success(meta).OkJSON()
}

func (h *Preview) DeleteOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")
	id := c.Param("id")
	if _, err := psvc.FindById(prodId); errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeBadRequest, "record: product with requested id not found").FailJSON()
	}

	err := svc.DeletePreview(id, prodId)
	if err != nil {
		return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(map[string]string{"id": id}).OkJSON()
}
