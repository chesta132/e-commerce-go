package errorlib

import (
	"errors"
	"fmt"
	"product-service/config"
	"strings"
)

var (
	ErrProductNotFound         = errors.New("record: product not found")
	ErrNoCategoryToCreate      = errors.New("record: no category to create product, create categories first before create product")
	ErrInvalidPreviewExt       = fmt.Errorf("validation: invalid preview extension, only allow %s", strings.Join(config.ALLOWED_PREVIEW_EXTENSION, ", "))
	ErrCantDeleteAllCategories = errors.New("conflict: can not delete all categories")
)
