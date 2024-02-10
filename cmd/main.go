package main

import (
	"flag"
	"iopairs-tool/internal"
	"os"
	"path/filepath"

	"log"

	"github.com/joho/godotenv"
)

var configPath string

func init() {
	// Parse the command-line flag for the config file path
	flag.StringVar(&configPath, "config", "config.yaml", "path to the config file")
}

func main() {
	flag.Parse()

	godotenv.Load()

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

	// Run the internal logic with the config and relative path
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
