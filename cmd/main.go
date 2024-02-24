package main

import (
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

	version = "undefined"
)

func init() {
	// Parse the command-line flag for the config file path
	flag.StringVar(&configPath, "config", "config.yaml", "path to the config file")
	flag.BoolVar(&showVersion, "version", false, "print the version")
}

func main() {
	flag.Parse()

	godotenv.Load()

	logo := `
 ________  ______   ______   ______   _________  
/_______/\/_____/\ /_____/\ /_____/\ /________/\ 
\__.::._\/\:::_ \ \\:::_ \ \\:::_ \ \\__.::.__\/ 
   \::\ \  \:\ \ \ \\:(🤖)\ \\:\ \ \ \  \::\ \   
   _\::\ \__\:\ \ \ \\: ___\/ \:\ \ \ \  \::\ \  
  /__\::\__/\\:\_\ \ \\ \ \    \:\/.:| |  \::\ \ 
  \________\/ \_____\/ \_\/     \____/_/   \__\/
  `
	fmt.Printf("\x1b[34m%s\x1b[0m", logo)

	if showVersion {
		fmt.Printf("\n	Version: %s\n", version)
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
