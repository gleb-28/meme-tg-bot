package utils

import "regexp"

func IsURL(text string) bool {
	re := regexp.MustCompile(`^(http|https)://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(text)
}
