package replylib

import (
	"net/http"

	"github.com/chesta132/goreply/reply"
)

const (
	CodeNotFound    = "NOT_FOUND"
	CodeServerError = "SERVER_ERROR"
	CodeBadRequest  = "BAD_REQUESt"
	CodeUnauthorized = "UNAUTHORIZED"
)

var Client = reply.NewClient(reply.Client{
	CodeAliases: map[string]int{
		CodeNotFound:    http.StatusNotFound,
		CodeServerError: http.StatusInternalServerError,
		CodeBadRequest: http.StatusBadRequest,
		CodeUnauthorized: http.StatusUnauthorized,
	},
})
