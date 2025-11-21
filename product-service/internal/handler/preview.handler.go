package handler

import (
	"bytes"
	"product-service/config"
	"product-service/internal/lib/errorlib"
	"product-service/internal/lib/previewlib"
	"product-service/internal/lib/replylib"
	"product-service/internal/service"

	"github.com/chesta132/e-commerce-go/shared/sreplylib"
	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
	"github.com/labstack/echo/v4"
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

	meta, preview, err := svc.GetPreviewFile(id, prodId)
	if err != nil {
		_, p, e := svc.ReadDefaultPreview()
		if e != nil {
			return rp.Error(sreplylib.CodeServerError, e.Error()).Info("Error while read default preview").FailJSON()
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
	if err != nil {
		return errorlib.HandleValidateProductOnPreviewError(err, rp)
	}

	return rp.Success(product.Previews).OkJSON()
}

func (h *Preview) CreateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")

	product, err := psvc.FindByIdWithRelation(prodId, []string{"Previews"})
	if err != nil {
		return errorlib.HandleValidateProductOnPreviewError(err, rp)
	}
	if len(product.Previews) >= config.MAX_PREVIEW {
		return rp.Error(sreplylib.CodeBadRequest, "record: max preview reached").FailJSON()
	}

	alt := c.FormValue("alt")
	fh, err := c.FormFile("file")
	if err != nil {
		return rp.Error(sreplylib.CodeBadRequest, err.Error()).FailJSON()
	}

	fbyte, err := svc.ResizePreview(fh)
	if err != nil {
		return rp.Error(sreplylib.CodeBadRequest, err.Error()).FailJSON()
	}

	meta := previewlib.GeneratePreview(fh, alt, prodId, true)
	err = svc.CreatePreviewWithFile(meta, fbyte)
	if err != nil {
		return rp.Error(sreplylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(meta).CreatedJSON()
}

func (h *Preview) UpdateOne(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	svc := h.svc.AttachEcho(c)
	psvc := h.psvc.AttachEcho(c)
	prodId := c.Param("prod-id")
	id := c.Param("id")

	_, err := psvc.FindById(prodId)
	if err != nil {
		return errorlib.HandleValidateProductOnPreviewError(err, rp)
	}

	meta, _, err := svc.GetPreviewFile(id, prodId)
	if err != nil {
		return errorlib.HandleQueryError(err, rp)
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return rp.Error(sreplylib.CodeBadRequest, err.Error()).FailJSON()
	}
	alt := c.FormValue("alt")
	meta.Alt = alt
	previewlib.MergePreview(fh, &meta)

	b, err := svc.ResizePreview(fh)
	if err != nil {
		return rp.Error(sreplylib.CodeBadRequest, err.Error()).FailJSON()
	}

	err = svc.UpdatePreviewWithFile(&meta, b)
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

	_, err := psvc.FindById(prodId)
	if err != nil {
		return errorlib.HandleValidateProductOnPreviewError(err, rp)
	}

	err = svc.DeletePreviewWithFile(id, prodId)
	if err != nil {
		return rp.Error(sreplylib.CodeServerError, err.Error()).FailJSON()
	}

	return rp.Success(map[string]string{"id": id}).OkJSON()
}
