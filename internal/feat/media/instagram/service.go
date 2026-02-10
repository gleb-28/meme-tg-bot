package instagram

import (
	"context"
	"memetgbot/internal/feat/media/video"
	"memetgbot/model"
)

type Extractor struct {
	video *video.Extractor
	image *ImageExtractor
}

func NewService(video *video.Extractor, image *ImageExtractor) *Extractor {
	return &Extractor{video: video, image: image}
}

func (e *Extractor) Extract(
	ctx context.Context,
	url string,
) (*model.MediaResult, error) {

	if res, err := e.video.Extract(ctx, url); err == nil {
		return res, nil
	}

	return e.image.Extract(ctx, url)
}
