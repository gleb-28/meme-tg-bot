package instagram

import (
	"context"
	"errors"
	"fmt"
	"memetgbot/internal/core/logger"
	"memetgbot/model"
	"memetgbot/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ImageConfig struct {
	GalleryDLPath  string
	UserAgent      string
	Sleep          string
	SleepRequest   string
	SleepExtractor string
	Sleep429       string
}

type ImageExtractor struct {
	downloadDir string
	cookiesPath string
	config      ImageConfig
	logger      logger.AppLogger
	sem         chan struct{}
	once        sync.Once
}

func NewImageService(downloadDir string, cookiesPath string, cfg ImageConfig, logger logger.AppLogger) *ImageExtractor {
	if cfg.GalleryDLPath == "" {
		cfg.GalleryDLPath = "gallery-dl"
	}

	return &ImageExtractor{
		downloadDir: downloadDir,
		cookiesPath: cookiesPath,
		config:      cfg,
		logger:      logger,
		sem:         make(chan struct{}, 1),
	}
}

func (e *ImageExtractor) Extract(
	ctx context.Context,
	url string,
) (*model.MediaResult, error) {
	if err := e.acquire(ctx); err != nil {
		return nil, err
	}
	defer e.release()

	if _, err := os.Stat(e.cookiesPath); err != nil {
		return nil, fmt.Errorf("instagram cookies file missing at %s: %w", e.cookiesPath, err)
	}

	ws, err := utils.NewTempWorkspace(e.downloadDir, "insta_")
	if err != nil {
		return nil, err
	}

	utils.RemoveAsyncDir(ws.Path, time.Minute)

	cmd := exec.CommandContext(
		ctx,
		e.config.GalleryDLPath,
		"--cookies", e.cookiesPath,
		"--user-agent", e.config.UserAgent,
		"--sleep", e.config.Sleep,
		"--sleep-request", e.config.SleepRequest,
		"--sleep-extractor", e.config.SleepExtractor,
		"--sleep-429", e.config.Sleep429,
		"-D", ws.Path,
		url,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, wrapGalleryDLError(err, string(out))
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

func (e *ImageExtractor) acquire(ctx context.Context) error {
	select {
	case e.sem <- struct{}{}:
		e.once.Do(func() {
			e.logger.Debug("instagram gallery-dl concurrency limit enabled: 1")
		})
		return nil
	case <-ctx.Done():
		return fmt.Errorf("instagram download queue wait canceled: %w", ctx.Err())
	}
}

func (e *ImageExtractor) release() {
	select {
	case <-e.sem:
	default:
	}
}

func wrapGalleryDLError(err error, output string) error {
	trimmed := strings.TrimSpace(output)
	lower := strings.ToLower(trimmed)

	if strings.Contains(lower, "we added a restriction to your account") ||
		strings.Contains(lower, "restriction to your account") {
		return fmt.Errorf(
			"instagram temporarily restricted this account. stop instagram downloads for this account, wait, refresh cookies from a real browser session, and keep slower gallery-dl pacing enabled. gallery-dl output: %s",
			trimmed,
		)
	}

	return fmt.Errorf("gallery-dl failed: %w; output: %s", err, trimmed)
}
