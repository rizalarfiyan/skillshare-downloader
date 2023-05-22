package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
)

func SafeName(filename string) string {
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

func ToSnakeCase(str string) string {
	tempStr := slug.MakeLang(str, "en")
	return strings.ReplaceAll(tempStr, "-", "_")
}

func MatchExtenstion(filename string, defaultExtension string) string {
	extension := defaultExtension
	ext := filepath.Ext(filename)
	if ext != "" {
		re := regexp.MustCompile(`\.(\w+)[\*\?]`)
		match := re.FindStringSubmatch(filename)
		if len(match) > 1 {
			extension = fmt.Sprintf(".%s", strings.ToLower(match[1]))
		}
	}
	return extension
}
