package previewlib

import (
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

func GetFilePath(preview *model.Preview, productId string) string {
	return filepath.Join(config.PREVIEW_PATH, productId, preview.ID+preview.Extension)
}

func GetDirPath(productId string) string {
	return filepath.Join(config.PREVIEW_PATH, productId)
}

func GetMetaPath(productId string) string {
	return filepath.Join(GetDirPath(productId), "meta.json")
}

func GeneratePreview(file *multipart.FileHeader, alt, productId string, withId bool) *model.Preview {
	id := ""
	if withId {
		id = uuid.NewString()
	}
	mime := ""
	h, ok := file.Header["Content-Type"]
	if ok && len(h) > 0 {
		mime = h[0]
	}
	if alt == "" {
		alt = GetFileName(file.Filename)
	}
	p := &model.Preview{
		ID: id,
		// Position:  len(meta.Preview) + 1, // fix auto increm by product id
		Extension: filepath.Ext(file.Filename),
		Alt:       alt,
		ProductId: productId,
		Mime:      mime,
	}
	p.Path = GetFilePath(p, productId)
	return p
}

func MergePreview(file *multipart.FileHeader, preview *model.Preview) {
	mime := ""
	h, ok := file.Header["Content-Type"]
	if ok && len(h) > 0 {
		mime = h[0]
	}

	preview.Extension = filepath.Ext(file.Filename)
	preview.Path = GetFilePath(preview, preview.ProductId)
	preview.Mime = mime
}

func GetDefaultPreview() model.Preview {
	path := config.DEFAULT_PREVIEW_PATH
	return model.Preview{
		ID:        path,
		Extension: filepath.Ext(path),
		Path:      path,
		Alt:       "default preview",
		Mime:      "image/png",
	}
}

func GetHeader(preview model.Preview, additional ...map[string]string) map[string]string {
	return map[string]string{
		"Content-Type":  preview.Mime,
		"X-Preview-Alt": preview.Alt,
	}
}

func ReadByHeader(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(f)
	f.Close()
	return b, err
}
