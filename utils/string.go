package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func SafeFolderName(filename string) string {
	str := DecodeAscii(filename)
	reSafe := regexp.MustCompile(`(?m)[^A-Za-z0-9-_ ]+`)
	reSpace := regexp.MustCompile(`\s+`)
	safeStr := reSafe.ReplaceAllString(str, "")
	return strings.TrimSpace(reSpace.ReplaceAllString(safeStr, " "))
}

func DecodeAscii(str string) string {
	r := regexp.MustCompile(`\\x[0-9A-Fa-f]{2}`)
	matches := r.FindAllString(str, -1)

	for _, match := range matches {
		hexStr := match[2:]
		hexNum, _ := strconv.ParseInt(hexStr, 16, 32)
		unicodeChar := rune(hexNum)
		str = strings.ReplaceAll(str, match, string(unicodeChar))
	}

	return str
}
