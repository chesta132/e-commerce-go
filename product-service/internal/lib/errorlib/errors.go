package errorlib

import "errors"

var (
	ErrProductNotFound    = errors.New("record: product not found")
	ErrNoCategoryToCreate = errors.New("record: no category to create product, create categories first before create product")
)
