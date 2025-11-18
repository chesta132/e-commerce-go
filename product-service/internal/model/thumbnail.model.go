package model

import "time"

type Thumbnail struct {
	ID        string    `json:"id"`
	Position  int       `json:"position"`
	Extension string    `json:"extension"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ThumbnailMeta struct {
	Thumbnails []*Thumbnail `json:"thumbnails"`
	ProjectId  string       `json:"projectId"`
}
