package utils

import (
	"errors"
	"fmt"
	"os"
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
