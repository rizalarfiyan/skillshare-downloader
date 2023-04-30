package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func GetCookieTxt(pathfile string) (string, error) {
	if !IsExistPath(pathfile) {
		return "", fmt.Errorf("%s not found", pathfile)
	}

	bytes, err := os.ReadFile(pathfile)
	if err != nil {
		return "", err
	}

	cookies := string(bytes)
	if cookies == "" {
		return "", errors.New("cookies is empty")
	}

	return cookies, nil
}

func CleanCookies(cookie string) string {
	re := regexp.MustCompile(`\n{2,}`)
	cookie = strings.TrimSpace(cookie)
	cookie = re.ReplaceAllString(cookie, "\n")
	return strings.ReplaceAll(cookie, "\n", " ")
}
