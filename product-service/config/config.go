package config

const (
	PAGINATION_LIMIT     = 100
	PREVIEW_PATH         = "/uploads/products/previews"
	DEFAULT_PREVIEW_PATH = "/uploads/products/previews/default_product_img.png"
	MAX_PREVIEW          = 15
	MAX_PREVIEW_SIZE     = "5M"
	MAX_IMAGE_WIDTH      = 480 // px
	MAX_VIDEO_WIDTH      = 480 // px
)

var (
	KEYWORD_QUERY                   = []string{"q", "query", "keyword"}
	ALLOWED_PREVIEW_IMAGE_EXTENSION = []string{".png", ".jpeg", ".jpg", ".webp"}
	ALLOWED_PREVIEW_VIDEO_EXTENSION = []string{".mp4", ".mov", ".webm"}
	ALLOWED_PREVIEW_EXTENSION       = append(ALLOWED_PREVIEW_IMAGE_EXTENSION, ALLOWED_PREVIEW_VIDEO_EXTENSION...)
)
