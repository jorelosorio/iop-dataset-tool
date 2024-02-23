package internal

import (
	"testing"
)

const textWithJSON = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit.
` + "```json\n" + `
{"conversations": [
	{"input": "What is the capital of France?", "output": "The capital of France is Paris."},
	{"input": "What is 2 + 2?", "output": "2 + 2 equals 4."}
]}
` + "```"

func TestFindJSONInText(t *testing.T) {
	responses, err := ExtractJsonFromText(textWithJSON)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if len(responses) != 1 {
		t.Errorf("Expected 1 response, got %d", len(responses))
	}
}
