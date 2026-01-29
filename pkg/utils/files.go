package utils

import "strings"

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
