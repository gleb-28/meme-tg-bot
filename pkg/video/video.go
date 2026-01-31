package video

import (
	"context"
	"fmt"
	"memetgbot/internal/core/logger"
	"memetgbot/pkg/utils"
	"os"
	"path/filepath"

	"github.com/lrstanley/go-ytdlp"
)

type Video struct {
	downloadDir string
	ytdlpPath   string
}

var VideoService = NewVideoService("./output", "/usr/bin/yt-dlp")

func NewVideoService(downloadDir string, ytdlpPath string) *Video {
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.MkdirAll(downloadDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("failed to create download directory %s: %v", downloadDir, err))
		}
	} else if err != nil {
		panic(fmt.Sprintf("failed to stat download directory %s: %v", downloadDir, err))
	}

	return &Video{
		downloadDir: downloadDir,
		ytdlpPath:   ytdlpPath,
	}
}

func (s *Video) DownloadVideo(ctx context.Context, videoURL string) (filePath string, fileName string, err error) {
	dl := ytdlp.New()

	if s.ytdlpPath != "" {
		dl = dl.SetExecutable(s.ytdlpPath)
	}

	dl = dl.
		Output(filepath.Join(s.downloadDir, "%(title)s.%(ext)s")).
		// Download the best video no better than 720p preferring framerate greater than 30, filesize < 50M
		// or the worst video (still preferring framerate greater than 30, filesize < 50M) if there is no such video,
		Format("((bv*[fps>30][filesize<50M]/bv*)[height<=720]/(wv*[fps>30][filesize<50M]/wv*)) + ba / (b[fps>30][filesize<50M]/b)[height<=480]/(w[fps>30]/w)").
		MergeOutputFormat("mp4").
		NoCheckCertificates()

	logger.Logger.Debug(fmt.Sprintf("Starting video download for: %s", videoURL))

	_, err = dl.Run(ctx, videoURL)
	if err != nil {
		return "", "", fmt.Errorf("download failed for %s: %w", videoURL, err)
	}

	logger.Logger.Debug(fmt.Sprintf("Video download completed for: %s. Scanning directory: %s", videoURL, s.downloadDir))

	entries, err := os.ReadDir(s.downloadDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to read download directory '%s': %w", s.downloadDir, err)
	}

	var latestFileInfo os.FileInfo
	var latestFilePath string

	for _, entry := range entries {
		info, err := entry.Info()
		name := utils.SanitizeFilename(entry.Name())
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to get file info for %s: %v", name, err))
			continue
		}

		if info.IsDir() {
			continue
		}

		if latestFileInfo == nil || info.ModTime().After(latestFileInfo.ModTime()) {
			latestFileInfo = info
			latestFilePath = filepath.Join(s.downloadDir, name)
		}
	}

	if latestFilePath == "" {
		return "", "", fmt.Errorf("no downloaded file found in '%s' after successful download command", s.downloadDir)
	}

	logger.Logger.Debug(fmt.Sprintf("Downloaded file found: %s", latestFilePath))
	return latestFilePath, utils.SanitizeFilename(latestFileInfo.Name()), nil
}
