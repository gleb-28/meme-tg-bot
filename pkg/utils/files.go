package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func SanitizeFilename(filename string) string {
	if len(filename) > 255 {
		filename = filename[:255]
	}

	invalidChars := []string{"/", "\\", "?", "%", "*", ":", "|", "\"", "<", ">"}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "_")
	}
	return filename
}

func RemoveCompressedSuffix(filePath string) string {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	base = strings.TrimSuffix(base, "_compressed")
	return base + ext
}

type FoundFile struct {
	Path string
	Name string
	Mod  time.Time
}

func FindImagesInDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".jpg" || ext == ".png" {
			safeName := SanitizeFilename(info.Name())
			safePath := filepath.Join(dir, safeName)
			if path != safePath {
				_ = os.Rename(path, safePath)
			}
			files = append(files, safePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, os.ErrNotExist
	}

	sort.Strings(files)
	return files, nil
}

type TempWorkspace struct {
	Path string
}

func NewTempWorkspace(baseDir, prefix string) (*TempWorkspace, error) {
	dir, err := os.MkdirTemp(baseDir, prefix)
	if err != nil {
		return nil, err
	}
	return &TempWorkspace{Path: dir}, nil
}

func RemoveAsync(path string, duration time.Duration) {
	go func(path string) {
		time.Sleep(duration)
		_ = os.Remove(path)
	}(path)
}

func RemoveAsyncDir(path string, duration time.Duration) {
	go func(path string) {
		time.Sleep(duration)
		_ = os.RemoveAll(path)
	}(path)
}

func ExtractLastStdoutLine(stdout string) string {
	lines := strings.Split(stdout, "\n")
	var res string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			res = line
			break
		}
	}
	return res
}
