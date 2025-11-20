package replylib

import (
	"net/http"
	"user-service/config"

	"github.com/chesta132/goreply/reply"
)

const (
	CodeNotFound     = "NOT_FOUND"
	CodeServerError  = "SERVER_ERROR"
	CodeBadRequest   = "BAD_REQUEST"
	CodeUnauthorized = "UNAUTHORIZED"
)

var CodeAliases = map[string]int{
	CodeNotFound:     http.StatusNotFound,
	CodeServerError:  http.StatusInternalServerError,
	CodeBadRequest:   http.StatusBadRequest,
	CodeUnauthorized: http.StatusUnauthorized,
}

var Client = reply.NewClient(reply.Client{
	CodeAliases: CodeAliases,
	DefaultHeaders: map[string]string{
		"X-Service": config.SERVICE,
	},
})

func GetCodeByStatus(status int) string {
	for k, v := range CodeAliases {
		if v == status {
			return k
		}
	}
	return CodeServerError
}
