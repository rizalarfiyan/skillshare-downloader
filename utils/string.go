package utils

import (
	"regexp"
	"strings"
)

func SafeFolderName(filename string) string {
	reSafe := regexp.MustCompile(`(?m)[^A-Za-z0-9-_ ]+`)
	reSpace := regexp.MustCompile(`\s+`)
	safeStr := reSafe.ReplaceAllString(filename, "")
	return strings.TrimSpace(reSpace.ReplaceAllString(safeStr, " "))
}
