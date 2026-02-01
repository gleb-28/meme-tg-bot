package video

import (
	"context"
	"fmt"
	"log"
	"memetgbot/internal/core/logger"
	"memetgbot/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

type VideoService struct {
	downloadDir string
	ytdlpPath   string
	cookiesPath string
	logger      logger.AppLogger
}

func MustNewVideoService(downloadDir string, ytdlpPath string, cookiesPath string, logger logger.AppLogger) *VideoService {
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.MkdirAll(downloadDir, 0755)
		if err != nil {
			log.Fatalf(fmt.Sprintf("failed to create download directory %s: %v", downloadDir, err))
		}
	} else if err != nil {
		log.Fatalf(fmt.Sprintf("failed to stat download directory %s: %v", downloadDir, err))
	}

	return &VideoService{
		downloadDir: downloadDir,
		ytdlpPath:   ytdlpPath,
		cookiesPath: cookiesPath,
		logger:      logger,
	}
}

func (videoService *VideoService) DownloadVideo(ctx context.Context, videoURL string) (filePath string, fileName string, err error) {
	dl := ytdlp.New()

	if videoService.ytdlpPath != "" {
		dl = dl.SetExecutable(videoService.ytdlpPath)
	}

	salt := utils.GenerateSalt()

	dl = dl.
		Output(filepath.Join(videoService.downloadDir, "%(title)s_"+salt+".%(ext)s")).
		// Download the best video no better than 720p preferring framerate greater than 30, filesize < 50M
		// or the worst video (still preferring framerate greater than 30, filesize < 50M) if there is no such video,
		Format("((bv*[fps>30][filesize<50M]/bv*)[height<=720]/(wv*[fps>30][filesize<50M]/wv*)) + ba / (b[fps>30][filesize<50M]/b)[height<=480]/(w[fps>30]/w)").
		MergeOutputFormat("mp4").
		NoCheckCertificates().
		Cookies(videoService.cookiesPath)

	videoService.logger.Debug(fmt.Sprintf("Starting video download for: %videoService", videoURL))

	_, err = dl.Run(ctx, videoURL)
	if err != nil {
		return "", "", fmt.Errorf("download failed for %videoService: %w", videoURL, err)
	}

	videoService.logger.Debug(fmt.Sprintf("VideoService download completed for: %videoService. Scanning directory: %videoService", videoURL, videoService.downloadDir))

	entries, err := os.ReadDir(videoService.downloadDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to read download directory '%videoService': %w", videoService.downloadDir, err)
	}

	var latestFilePath string
	var latestFileName string

	for _, entry := range entries {
		info, err := entry.Info()
		latestFileName = utils.SanitizeFilename(entry.Name())
		if err != nil {
			videoService.logger.Error(fmt.Sprintf("Failed to get file info for %videoService: %v", latestFileName, err))
			continue
		}

		if info.IsDir() {
			continue
		}

		if strings.Contains(latestFileName, salt) {
			latestFilePath = filepath.Join(videoService.downloadDir, latestFileName)
			break
		}
	}

	if latestFilePath == "" {
		return "", "", fmt.Errorf("no downloaded file found in '%videoService' after successful download command", videoService.downloadDir)
	}

	videoService.logger.Debug(fmt.Sprintf("Downloaded file found: %videoService", latestFilePath))
	return latestFilePath, latestFileName, nil
}

func (videoService *VideoService) DeleteVideoByName(videoName string) error {
	videoPath := filepath.Join(videoService.downloadDir, utils.SanitizeFilename(videoName))

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return fmt.Errorf("video %videoService does not exist in %videoService", videoName, videoService.downloadDir)
	} else if err != nil {
		return fmt.Errorf("failed to stat video file %videoService: %w", videoPath, err)
	}

	if err := os.Remove(videoPath); err != nil {
		return fmt.Errorf("failed to delete video %videoService: %w", videoName, err)
	}

	videoService.logger.Debug(fmt.Sprintf("VideoService %videoService has been deleted successfully", videoName))
	return nil
}
