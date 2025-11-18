package errorlib

import "errors"

var (
	ErrProductNotFound       = errors.New("record: product not found")
	ErrNoCategoryToCreate    = errors.New("record: no category to create product, create categories first before create product")
	ErrThumbnailNotFound     = errors.New("record: thumbnail not found")
	ErrThumbnailMetaNotFound = errors.New("record: thumbnail meta not found")
)
