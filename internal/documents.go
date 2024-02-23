package internal

import (
	"os"
	"path/filepath"
)

func GetDocumentsCorpus(documents []string, relativePath string) (string, error) {
	documentPaths, err := getDocumentPaths(documents, relativePath)
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

func getDocumentPaths(documents []string, relativePath string) ([]string, error) {
	uniqRelativeDocumentPaths := map[string]string{}
	for _, documentPath := range documents {
		relativeDocumentPath := filepath.Join(relativePath, documentPath)
		relativeDocumentPaths, err := ExpandFiles(relativeDocumentPath)
		if err != nil {
			return []string{}, err
		}

		for _, relativeDocumentPath := range relativeDocumentPaths {
			uniqRelativeDocumentPaths[relativeDocumentPath] = relativeDocumentPath
		}
	}

	var documentPaths []string
	for _, relativeDocumentPath := range uniqRelativeDocumentPaths {
		documentPaths = append(documentPaths, relativeDocumentPath)
	}

	return documentPaths, nil
}
