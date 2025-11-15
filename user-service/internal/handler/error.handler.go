package handler

import (
	"fmt"
	"user-service/internal/lib/replylib"

	adapter "github.com/chesta132/goreply/adapter/echo"
	"github.com/labstack/echo/v4"
)

func NotFound(c echo.Context) error {
	rp := replylib.Client.New(adapter.AdaptEcho(c))
	return rp.Error(replylib.CodeNotFound, fmt.Sprintf("endpoint: path %s not found", c.Request().URL)).FailJSON()
}
