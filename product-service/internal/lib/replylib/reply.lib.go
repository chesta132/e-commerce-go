package replylib

import (
	"net/http"
	"product-service/config"

	"github.com/chesta132/goreply/reply"
)

const (
	CodeNotFound     = "NOT_FOUND"
	CodeServerError  = "SERVER_ERROR"
	CodeBadRequest   = "BAD_REQUEST"
	CodeBadGateway   = "BAD_GATEWAY"
	CodeUnauthorized = "UNAUTHORIZED"
)

var Client = reply.NewClient(reply.Client{
	CodeAliases: map[string]int{
		CodeNotFound:     http.StatusNotFound,
		CodeServerError:  http.StatusInternalServerError,
		CodeBadRequest:   http.StatusBadRequest,
		CodeUnauthorized: http.StatusUnauthorized,
		CodeBadGateway: http.StatusBadGateway,
	},
	DefaultHeaders: map[string]string{
		"X-Service": config.SERVICE,
	},
})
