package thumbnaillib

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"product-service/config"
	"product-service/internal/model"
	"sort"
	"strings"

	"github.com/google/uuid"
)

func Reindex(meta *model.ThumbnailMeta) {
	sort.Slice(meta.Thumbnails, func(i, j int) bool {
		return meta.Thumbnails[i].Position < meta.Thumbnails[j].Position
	})
	for i := range meta.Thumbnails {
		meta.Thumbnails[i].Position = i + 1
	}
}

func GetExtension(fileName string) string {
	i := strings.LastIndex(fileName, ".")
	return fileName[i+1:]
}

func GetFilePath(thumbnail *model.Thumbnail, projectId string) string {
	return filepath.Join(config.THUMBNAIL_PATH, projectId, fmt.Sprintf("%s.%s", thumbnail.ID, thumbnail.Extension))
}

func GetDirPath(projectId string) string {
	return filepath.Join(config.THUMBNAIL_PATH, projectId)
}

func GetMetaPath(projectId string) string {
	return filepath.Join(GetDirPath(projectId), "meta.json")
}

func GenerateThumbnail(file *multipart.FileHeader, meta *model.ThumbnailMeta, withId bool) *model.Thumbnail {
	id := ""
	if withId {
		id = uuid.NewString()
	}
	t := &model.Thumbnail{
		ID:        id,
		Position:  len(meta.Thumbnails) + 1,
		Extension: GetExtension(file.Filename),
	}
	t.Path = GetFilePath(t, meta.ProjectId)
	return t
}
