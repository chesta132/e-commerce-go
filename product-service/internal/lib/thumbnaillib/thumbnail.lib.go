package thumbnaillib

import (
	"fmt"
	"path/filepath"
	"product-service/config"
	"product-service/internal/model"
	"sort"
)

func Reindex(meta *model.ThumbnailMeta) {
	sort.Slice(meta.Thumbnails, func(i, j int) bool {
		return meta.Thumbnails[i].Position < meta.Thumbnails[j].Position
	})
	for i := range meta.Thumbnails {
		meta.Thumbnails[i].Position = i + 1
	}
}

func GetFilePath(thumbnail *model.Thumbnail, projectId string) string {
	return filepath.Join(config.THUMBNAIL_PATH, projectId, fmt.Sprintf("%s.%s", thumbnail.ID, thumbnail.Extension))
}

func GetDirPath(projectId string) string {
	return filepath.Join(config.THUMBNAIL_PATH, projectId)
}
