package errorlib

import (
	"errors"
	"strings"
	"user-service/internal/lib/replylib"

	"github.com/chesta132/goreply/reply"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func HandleSigninError(err error, rp *reply.Reply) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rp.Error(replylib.CodeNotFound, ErrUserNotFound.Error()).FailJSON()
		return
	} else if errors.Is(err, ErrWrongPassword) {
		rp.Error(replylib.CodeBadRequest, err.Error(), reply.OptErrorPayload{Field: "password"}).FailJSON()
		return
	} else if e, ok := err.(validator.ValidationErrors); ok {
		fields := []string{}
		for _, fe := range e {
			fields = append(fields, fe.Field())
		}
		rp.Error(replylib.CodeBadRequest, e.Error(), reply.OptErrorPayload{Field: strings.Join(fields, ", ")}).FailJSON()
		return
	}
	rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}
