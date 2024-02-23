package internal

import (
	"strings"
	"unicode"
)

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
