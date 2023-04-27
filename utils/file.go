package utils

import "os"

func IsExistPath(pathname string) bool {
	_, err := os.Stat(pathname)
	return !os.IsNotExist(err)
}

func CreateDir(pathname string) error {
	if !IsExistPath(pathname) {
		err := os.MkdirAll(pathname, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}
