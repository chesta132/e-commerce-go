package errorlib

import (
	"errors"
	"product-service/internal/lib/replylib"

	"github.com/chesta132/goreply/reply"
	"gorm.io/gorm"
)

func HandleQueryError(err error, rp *reply.Reply) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeNotFound, err.Error()).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}
