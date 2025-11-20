package errorlib

import (
	"errors"
	"product-service/internal/lib/replylib"
	"product-service/internal/lib/thumbnaillib"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func HandleQueryError(err error, rp *reply.Reply) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeNotFound, err.Error()).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}

func HandleValidateError(err validator.ValidationErrors, rp *reply.Reply) error {
	fields := []string{}
	for _, fe := range err {
		fields = append(fields, fe.Field())
	}
	return rp.Error(replylib.CodeBadRequest, err.Error(), reply.OptErrorPayload{Field: strings.Join(fields, ", ")}).FailJSON()
}

func HandleCreateProductError(err error, rp *reply.Reply) error {
	if err, ok := err.(validator.ValidationErrors); ok {
		return HandleValidateError(err, rp)
	}
	if errors.Is(err, ErrNoCategoryToCreate) {
		return rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}

func HandleGetThumbnailError(err error, defaultThumb []byte, rp *reply.Reply) error {
	errHeader := map[string]string{"X-Error": err.Error()}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errHeader["X-Error"] = "record: thumbnail not found, fallback to default thumbnail"
	}
	meta := thumbnaillib.GetDefaultThumbnail()
	return rp.Success(defaultThumb).AddHeaders(thumbnaillib.GetHeader(meta, errHeader)).ReplyBinary(404)
}
