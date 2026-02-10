package video

import (
	"context"
	"fmt"
	"log"
	"memetgbot/internal/core/logger"
	"memetgbot/model"
	"memetgbot/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

type Extractor struct {
	downloadDir string
	ytdlpPath   string
	ffmpegPath  string
	cookiesPath string
	logger      logger.AppLogger
}

func MustNewVideoService(downloadDir string, ytdlpPath string, cookiesPath string, ffmpegPath string, logger logger.AppLogger) *Extractor {
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.MkdirAll(downloadDir, 0755)
		if err != nil {
			log.Fatalf("failed to create download directory %s: %v", downloadDir, err)
		}
	} else if err != nil {
		log.Fatalf("failed to stat download directory %s: %v", downloadDir, err)
	}

	return &Extractor{
		downloadDir: downloadDir,
		ytdlpPath:   ytdlpPath,
		cookiesPath: cookiesPath,
		ffmpegPath:  ffmpegPath,
		logger:      logger,
	}
}

func (videoService *Extractor) Extract(ctx context.Context, videoURL string) (*model.MediaResult, error) {
	dl := ytdlp.New()

	if videoService.ytdlpPath != "" {
		dl = dl.SetExecutable(videoService.ytdlpPath)
	}

	dl = dl.
		Output(filepath.Join(videoService.downloadDir, "%(title)s.%(ext)s")).
		// Download the best video no better than 720p preferring framerate greater than 30, filesize < 50M
		// or the worst video (still preferring framerate greater than 30, filesize < 50M) if there is no such video,
		Format("((bv*[fps>30][filesize<50M]/bv*)[height<=720]/(wv*[fps>30][filesize<50M]/wv*)) + ba / (b[fps>30][filesize<50M]/b)[height<=480]/(w[fps>30]/w)").
		MergeOutputFormat("mp4").
		NoCheckCertificates().
		Cookies(videoService.cookiesPath).
		Print("after_move:filepath")

	videoService.logger.Debug(fmt.Sprintf("Starting video download for: %v", videoURL))

	res, err := dl.Run(ctx, videoURL)
	if err != nil {
		return nil, fmt.Errorf("download failed for %v: %w", videoURL, err)
	}

	if len(res.Stdout) == 0 {
		return nil, fmt.Errorf("yt-dlp returned no output file")
	}

	videoService.logger.Debug(fmt.Sprintf("video download completed for: %v. Scanning directory: %v", videoURL, videoService.downloadDir))

	downloadedPath := utils.ExtractLastStdoutLine(res.Stdout)

	if downloadedPath == "" {
		return nil, fmt.Errorf("yt-dlp returned no output filepath")
	}

	fileName := filepath.Base(downloadedPath)

	compressedPath, err := videoService.compressVideo(downloadedPath)
	if err != nil {
		videoService.logger.Error(err.Error())
		utils.RemoveAsync(downloadedPath, time.Minute)
		return &model.MediaResult{
			Type: model.MediaVideo,
			Files: []model.MediaFile{
				{Path: downloadedPath, Name: fileName},
			},
		}, nil
	}

	_ = os.Remove(downloadedPath)
	utils.RemoveAsync(compressedPath, time.Minute)

	return &model.MediaResult{
		Type: model.MediaVideo,
		Files: []model.MediaFile{
			{Path: compressedPath, Name: filepath.Base(compressedPath)},
		},
	}, nil
}

func (videoService *Extractor) compressVideo(inputPath string) (string, error) {
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_compressed.mp4"
	videoService.logger.Debug(fmt.Sprintf("Compressing video: %s", inputPath))

	cmd := exec.Command(
		videoService.ffmpegPath,
		"-i", inputPath,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "25",
		"-c:a", "aac",
		"-y",
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("compression failed: %w", err)
	}

	if err := os.Remove(inputPath); err != nil {
		videoService.logger.Error(fmt.Sprintf("Failed to remove original video: %s", err))
	}

	videoService.logger.Debug(fmt.Sprintf("Compressed video created: %s", outputPath))
	return outputPath, nil
}
