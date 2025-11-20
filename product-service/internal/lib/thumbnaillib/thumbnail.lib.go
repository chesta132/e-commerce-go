package thumbnaillib

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"product-service/config"
	"product-service/internal/model"
	"strings"

	"github.com/google/uuid"
)

// returns "" if file name is not a valid file name
func GetFileName(fileName string) string {
	if fn := filepath.Ext(fileName); fn == "" {
		return ""
	}
	indexs := []int{strings.LastIndex(fileName, "."), strings.LastIndex(fileName, "/"), strings.LastIndex(fileName, "\\")}
	for _, i := range indexs {
		if i != -1 {
			fileName = fileName[:i]
		}
	}
	return fileName
}

func GetFilePath(thumbnail *model.Thumbnail, productId string) string {
	return filepath.Join(config.THUMBNAIL_PATH, productId, fmt.Sprintf("%s.%s", thumbnail.ID, thumbnail.Extension))
}

func GetDirPath(productId string) string {
	return filepath.Join(config.THUMBNAIL_PATH, productId)
}

func GetMetaPath(productId string) string {
	return filepath.Join(GetDirPath(productId), "meta.json")
}

func GenerateThumbnail(file *multipart.FileHeader, alt, productId string, withId bool) *model.Thumbnail {
	id := ""
	if withId {
		id = uuid.NewString()
	}
	mime := ""
	h, ok := file.Header["Content-Type"]
	if ok && len(h) > 0 {
		mime = h[0]
	}
	t := &model.Thumbnail{
		ID: id,
		// Position:  len(meta.Thumbnails) + 1, // fix auto increm by product id
		Extension: filepath.Ext(file.Filename),
		Alt:       alt,
		ProductId: productId,
		Mime:      mime,
	}
	t.Path = GetFilePath(t, productId)
	return t
}

func MergeThumbnail(file *multipart.FileHeader, thumbnail *model.Thumbnail) {
	mime := ""
	h, ok := file.Header["Content-Type"]
	if ok && len(h) > 0 {
		mime = h[0]
	}

	thumbnail.Extension = filepath.Ext(file.Filename)
	thumbnail.Path = GetFilePath(thumbnail, thumbnail.ProductId)
	thumbnail.Mime = mime
}

func GetDefaultThumbnail() model.Thumbnail {
	path := config.DEFAULT_THUMBNAIL_PATH
	return model.Thumbnail{
		ID:        path,
		Extension: filepath.Ext(path),
		Path:      path,
		Alt:       "default thumbnail",
		Mime:      "image/png",
	}
}

func GetHeader(thumbnail model.Thumbnail, additional ...map[string]string) map[string]string {
	return map[string]string{
		"Content-Type":    thumbnail.Mime,
		"X-Thumbnail-Alt": thumbnail.Alt,
	}
}

func ReadByHeader(fh *multipart.FileHeader) (multipart.File, []byte, error) {
	f, err := fh.Open()
	if err != nil {
		return f, nil, err
	}
	b, err := io.ReadAll(f)
	return f, b, err
}
