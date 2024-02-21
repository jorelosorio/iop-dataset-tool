package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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

func ExtractFileName(path string) (string, error) {
	if IsFile(path) {
		return filepath.Base(path), nil
	}

	return "", fmt.Errorf("path is not a file")
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

func SplitTextIntoChunks(text string, chunkSize int) []string {
	var splits []string
	var chunk strings.Builder

	runes := []rune(text)
	for startIdx := 0; startIdx < len(runes); {
		endIdx := min(startIdx+chunkSize, len(runes))
		if endIdx < len(runes) {
			// Extend to complete the word if the word is being split
			for endIdx < len(runes) && !unicode.IsSpace(runes[endIdx]) && !unicode.IsPunct(runes[endIdx]) {
				endIdx++
			}
		}

		chunk.WriteString(string(runes[startIdx:endIdx]))
		splits = append(splits, strings.TrimSpace(chunk.String()))
		chunk.Reset()

		// Start the next chunk right after the current one ends
		startIdx = endIdx
	}

	return splits
}
