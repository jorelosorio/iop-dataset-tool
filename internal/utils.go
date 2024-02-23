package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func FolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func MakeDir(path string) error {
	if !FolderExists(path) {
		return os.MkdirAll(path, 0755)
	}

	return nil
}

// Is a file name
func IsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !fileInfo.IsDir()
}

func ExpandFiles(file string) ([]string, error) {
	matches, err := filepath.Glob(file)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files found for: %s", file)
	}

	return matches, nil
}
