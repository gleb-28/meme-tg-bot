package instagram

import (
	"context"
	"errors"
	"memetgbot/internal/core/logger"
	"memetgbot/model"
	"memetgbot/pkg/utils"
	"os/exec"
	"path/filepath"
	"time"
)

type ImageExtractor struct {
	downloadDir string
	cookiesPath string
	logger      logger.AppLogger
}

func NewImageService(downloadDir string, cookiesPath string, logger logger.AppLogger) *ImageExtractor {
	return &ImageExtractor{downloadDir: downloadDir, cookiesPath: cookiesPath, logger: logger}
}

func (e *ImageExtractor) Extract(
	ctx context.Context,
	url string,
) (*model.MediaResult, error) {

	ws, err := utils.NewTempWorkspace(e.downloadDir, "insta_")
	if err != nil {
		return nil, err
	}

	utils.RemoveAsyncDir(ws.Path, time.Minute)

	cmd := exec.CommandContext(
		ctx,
		"gallery-dl",
		"--cookies", e.cookiesPath,
		"-D", ws.Path,
		url,
	)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	e.logger.Debug("gallery-dl instagram image download")

	files, err := utils.FindImagesInDir(ws.Path)
	if err != nil {
		return nil, errors.New("no images downloaded")
	}

	mediaFiles := make([]model.MediaFile, len(files))
	for i, path := range files {
		mediaFiles[i] = model.MediaFile{Path: path, Name: filepath.Base(path)}
	}

	return &model.MediaResult{
		Type:  model.MediaAlbum,
		Files: mediaFiles,
	}, nil
}
