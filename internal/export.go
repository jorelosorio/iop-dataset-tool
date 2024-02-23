package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func ExportData(data any, outputDir string) error {
	if dataString, isString := data.(string); isString {
		err := dumpDataIntoOutputDir([]byte(dataString), "txt", outputDir)
		if err != nil {
			return err
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = dumpDataIntoOutputDir(jsonData, "jsonl", outputDir)
	if err != nil {
		return err
	}

	return nil
}

func dumpDataIntoOutputDir(data []byte, extension, outputDir string) error {
	err := MakeDir(outputDir)
	if err != nil {
		return err
	}

	var fileName string
	currentTime := time.Now().Unix()
	fileName = fmt.Sprintf("%d.%s", currentTime, extension)
	filePath := filepath.Join(outputDir, fileName)

	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
