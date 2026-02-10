package model

import "context"

type MediaResult struct {
	Type  MediaType
	Files []MediaFile
}

type MediaType string

const (
	MediaVideo MediaType = "video"
	MediaAlbum MediaType = "album"
)

type MediaFile struct {
	Path string
	Name string
}

type MediaExtractor interface {
	Extract(ctx context.Context, url string) (*MediaResult, error)
}
