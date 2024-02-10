package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func exportProcessResponseData(rawJSONResponse Response, outputDir string) error {
	jsonData, err := json.MarshalIndent(rawJSONResponse, "", "  ")
	if err != nil {
		return err
	}

	error := DumpDataIntoOutputDir(jsonData, outputDir, "")
	if error != nil {
		return error
	}

	return nil
}

func DumpDataIntoOutputDir(data []byte, outputDir string, sufix string) error {
	err := MakeDir(outputDir)
	if err != nil {
		return err
	}

	var fileName string
	currentTime := time.Now().Unix()
	if sufix == "" {
		fileName = fmt.Sprintf("%d.json", currentTime)
	} else {
		fileName = fmt.Sprintf("%d_%s.json", currentTime, sufix)
	}

	filePath := filepath.Join(outputDir, fileName)

	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
