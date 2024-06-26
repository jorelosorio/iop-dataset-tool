package main

import (
	_ "embed"
	"flag"
	"fmt"
	"iopairs-tool/internal"
	"os"
	"path/filepath"

	"log"

	"github.com/joho/godotenv"
)

var (
	configPath  string
	showVersion bool

	//go:embed logo.txt
	asciiLogo string
)

var (
	version = "dev"
)

func init() {
	// Parse the command-line flag for the config file path
	flag.StringVar(&configPath, "config", "config.yaml", "path to the config file")
	flag.BoolVar(&showVersion, "version", false, "print the version")
}

func main() {
	flag.Parse()

	godotenv.Load()

	fmt.Println(internal.PsychedelicGradient(asciiLogo))

	if showVersion {
		fmt.Println(internal.Colorizef(internal.ColorMagenta, "	\n	Version: %s", version))
		os.Exit(0)
	}

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
