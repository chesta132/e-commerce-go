package errorlib

import (
	"errors"
	"strings"
	"user-service/internal/lib/replylib"

	"github.com/chesta132/goreply/reply"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func HandleSigninError(err error, rp *reply.Reply) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeNotFound, ErrUserNotFound.Error()).FailJSON()
	} else if errors.Is(err, ErrWrongPassword) {
		return rp.Error(replylib.CodeBadRequest, err.Error(), reply.OptErrorPayload{Field: "password"}).FailJSON()
	} else if e, ok := err.(validator.ValidationErrors); ok {
		fields := []string{}
		for _, fe := range e {
			fields = append(fields, fe.Field())
		}
		return rp.Error(replylib.CodeBadRequest, e.Error(), reply.OptErrorPayload{Field: strings.Join(fields, ", ")}).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}

func HandleSignupError(err error, rp *reply.Reply) error {
	if errors.Is(err, ErrEmailRegistered) {
		return rp.Error(replylib.CodeBadRequest, err.Error(), reply.OptErrorPayload{Field: "email"}).FailJSON()
	} else if e, ok := err.(validator.ValidationErrors); ok {
		fields := []string{}
		for _, fe := range e {
			fields = append(fields, fe.Field())
		}
		return rp.Error(replylib.CodeBadRequest, e.Error(), reply.OptErrorPayload{Field: strings.Join(fields, ", ")}).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}

func HandleValidateAuthError(err error, rp *reply.Reply) error {
	if errors.Is(err, ErrTokenExpired) || errors.Is(err, ErrInvalidToken) {
		return rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}

func HandleQueryError(err error, rp *reply.Reply) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Error(replylib.CodeNotFound, err.Error()).FailJSON()
	}
	return rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
}
