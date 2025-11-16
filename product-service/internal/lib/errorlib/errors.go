package errorlib

import "errors"

var (
	ErrProductNotFound = errors.New("record: product not found")
)
