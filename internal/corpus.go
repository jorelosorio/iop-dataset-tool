package internal

import "os"

func CollectCorpusData(files []string, relativePath string) (string, error) {
	documentPaths, err := getDocumentPaths(files, relativePath)
	if err != nil {
		return "", err
	}

	var corpus string
	for _, documentPath := range documentPaths {
		text, err := os.ReadFile(documentPath)
		if err != nil {
			return "", err
		}

		corpus += string(text) + "\n"
	}

	return corpus, nil
}
