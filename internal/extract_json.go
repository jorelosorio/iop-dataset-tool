package internal

import (
	"fmt"
	"regexp"
	"strings"
)

func ExtractJsonFromText(text string) ([]string, error) {
	// Define a regular expression to find ```json ... ```
	r, err := regexp.Compile("```json([\\s\\S]*?)```")
	if err != nil {
		return nil, err
	}

	matches := r.FindAllStringSubmatch(text, -1)
	if matches == nil {
		return nil, fmt.Errorf("no JSON found in response")
	}

	var results []string
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		jsonStr := match[1]
		jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

		results = append(results, jsonStr)
	}

	return results, nil
}
