package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
	reSafe := regexp.MustCompile(`(?m)[^A-Za-z0-9- ]+`)
	safeStr := reSafe.ReplaceAllString(str, "")
	str = strings.ReplaceAll(safeStr, " ", "_")
	var sb strings.Builder
	var prev rune
	for _, curr := range str {
		if curr >= 'A' && curr <= 'Z' {
			if prev >= 'a' && prev <= 'z' {
				sb.WriteRune('_')
			}
			sb.WriteRune(curr + ('a' - 'A'))
		} else {
			sb.WriteRune(curr)
		}
		prev = curr
	}
	return sb.String()
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
