package utils

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
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

func GenerateSalt() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate salt")
	}
	return hex.EncodeToString(b)
}

func RemoveSaltFromFileName(fileName string) string {
	ext := filepath.Ext(fileName)

	baseName := strings.TrimSuffix(fileName, ext)

	lastUnderscore := strings.LastIndex(baseName, "_")
	if lastUnderscore == -1 {
		return fileName // если подчеркивания нет, возвращаем оригинальное имя
	}

	return baseName[:lastUnderscore] + ext
}

func RemoveCompressedSuffix(filePath string) string {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	base = strings.TrimSuffix(base, "_compressed")
	return base + ext
}
