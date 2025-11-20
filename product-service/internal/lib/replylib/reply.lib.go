package replylib

import (
	"product-service/config"

	"github.com/chesta132/e-commerce-go/shared/sreplylib"
	"github.com/chesta132/goreply/reply"
)

var Client = reply.NewClient(reply.Client{
	CodeAliases: sreplylib.CodeAliases,
	DefaultHeaders: map[string]string{
		"X-Service": config.SERVICE,
	},
})
