package handler

import (
	"fmt"
	"net/http"
	"product-service/internal/lib/replylib"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/chesta132/goreply/reply"
	"github.com/labstack/echo/v4"
)

func NotFound(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	return rp.Error(replylib.CodeNotFound, fmt.Sprintf("endpoint: path %s not found", c.Request().URL)).FailJSON()
}

func DefaultErrorHandler(err error, c echo.Context) {
	rp := replylib.Client.New(adapter.AdaptEcho(c))

	he, ok := err.(*echo.HTTPError)
	if !ok {
		he = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	details := ""
	if he.Internal != nil {
		details = he.Internal.Error()
	}

	code := replylib.GetCodeByStatus(he.Code)
	rp.Error(code, fmt.Sprint(he.Message), reply.OptErrorPayload{
		Details: details,
	}).FailJSON(he.Code)
}
