package media

import (
	"context"
	"memetgbot/internal/feat/media/instagram"
	"memetgbot/internal/feat/media/video"
	"memetgbot/model"
	"strings"
)

type Service struct {
	insta *instagram.Extractor
	video *video.Extractor
}

func NewService(insta *instagram.Extractor, video *video.Extractor) *Service {
	return &Service{insta: insta, video: video}
}

func (s *Service) Extract(
	ctx context.Context,
	url string,
) (*model.MediaResult, error) {

	if strings.Contains(url, "instagram.com") {
		return s.insta.Extract(ctx, url)
	}

	return s.video.Extract(ctx, url)
}
