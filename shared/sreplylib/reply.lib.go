package sreplylib

import "net/http"

const (
	CodeNotFound     = "NOT_FOUND"
	CodeServerError  = "SERVER_ERROR"
	CodeBadRequest   = "BAD_REQUEST"
	CodeBadGateway   = "BAD_GATEWAY"
	CodeUnauthorized = "UNAUTHORIZED"
)

var CodeAliases = map[string]int{
	CodeNotFound:     http.StatusNotFound,
	CodeServerError:  http.StatusInternalServerError,
	CodeBadRequest:   http.StatusBadRequest,
	CodeUnauthorized: http.StatusUnauthorized,
	CodeBadGateway:   http.StatusBadGateway,
}

func GetCodeByStatus(status int) string {
	for k, v := range CodeAliases {
		if v == status {
			return k
		}
	}
	return CodeServerError
}
