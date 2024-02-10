package main

import (
	"iopairs-tool/internal"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	configPath := "../data/quillero/config.yaml"
	config, err := internal.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Current directory
	// Find the relative path to the desired directory
	relativePath, err := getConfigRelativePath(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = internal.Run(config, relativePath)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigRelativePath(configPath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, filepath.Dir(configPath)), nil
}
