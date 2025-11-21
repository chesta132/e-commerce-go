package previewlib

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"product-service/config"
	"product-service/internal/model"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func ResizeImage(src multipart.File) ([]byte, error) {
	img, _ := imaging.Decode(src)
	resized := imaging.Resize(img, config.MAX_IMAGE_WIDTH, 0, imaging.Lanczos)

	out := new(bytes.Buffer)
	imaging.Encode(out, resized, imaging.JPEG)

	return out.Bytes(), nil
}

func ResizeVideo(src multipart.File, ext string) ([]byte, error) {
	tmp, err := os.CreateTemp("", "upload-*"+ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, src)
	if err != nil {
		return nil, err
	}
	tmp.Close()

	outPath := "/tmp/ffmpeg-" + uuid.NewString() + ".mp4"

	err = ffmpeg.Input(tmp.Name()).
		Output(outPath, ffmpeg.KwArgs{
			"vf":  fmt.Sprintf("scale=%d:trunc(ow/a/2)*2", config.MAX_VIDEO_WIDTH),
			"c:v": "libx264",
		}).
		OverWriteOutput().
		WithErrorOutput(os.Stderr).
		Run()

	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w", err)
	}

	out, err := os.ReadFile(outPath)
	if err != nil {
		return nil, err
	}
	os.Remove(outPath)

	return out, nil
}
