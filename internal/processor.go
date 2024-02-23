package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/xeipuuv/gojsonschema"
)

func Run(config Config, relativePath string) error {
	for _, process := range config.Processes {
		if process.Skip {
			log.Printf("🚫 Skipping: %s", process.Name)
			continue
		}

		outputDir := filepath.Join(relativePath, process.OutputDir)

		log.Printf("🤓 Processing: %s", process.Name)

		log.Printf("📚 Collecting corpus...")
		corpus, err := GetDocumentsCorpus(process.Documents, relativePath)
		if err != nil {
			return err
		}

		log.Printf("📑 Corpus collected: %d", len(corpus))

		if len(corpus) == 0 {
			log.Printf("🚫 Skipping: %s, empty corpus", process.Name)
			continue
		}

		chunkCorpus := SplitTextIntoChunks(corpus, process.ChunkSize)
		log.Printf("🍤🍤 Chunking corpus into %d", len(chunkCorpus))
		for chunkIndex, corpus := range chunkCorpus {
			chunkInfo := fmt.Sprintf("🍤 Chunk (%d of %d)", chunkIndex+1, len(chunkCorpus))

			log.Printf("🥁 %s processing", chunkInfo)

			target, err := config.GetTarget(process.Target)
			if err != nil {
				return err
			}

			data, err := requestModel(
				string(corpus),
				process,
				target,
			)

			// It might that the JSON is not formatted correctly, or an error was found
			// then save raw JSON to a file for debugging
			if err != nil {
				log.Printf("	💥 An error occurred: %v", err)
				log.Print("	🛟 Saving response to a file for debugging")

				errExporting := ExportData(data, outputDir)
				if errExporting != nil {
					return errExporting
				}

				return err
			}

			err = ExportData(data, outputDir)
			if err != nil {
				return err
			}

			log.Printf("	💾 JSON response saved")
			log.Printf("	✅ %s processing completed!", chunkInfo)
		}
	}

	log.Printf("🎉🍻 All processes completed successfully!")

	return nil
}

func requestModel(corpus string, process Process, target Target) (any, error) {
	templateData := struct {
		Document string
	}{Document: string(corpus)}

	systemPrompt, err := parseTemplate("systemPrompt", process.SystemPrompt, templateData)
	if err != nil {
		return nil, fmt.Errorf("error executing system prompt template for %s: %v", process.Name, err)
	}

	userPrompt, err := parseTemplate("userPrompt", process.UserPrompt, templateData)
	if err != nil {
		return nil, fmt.Errorf("error executing user prompt template for %s: %v", process.Name, err)
	}

	req := openai.ChatCompletionRequest{
		Model: process.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
	}

	if process.MaxTokens != 0 {
		req.MaxTokens = process.MaxTokens
	}

	if process.Temperature != 0 {
		req.Temperature = process.Temperature
	}

	apiToken := os.Getenv(target.ApiKeyEnv)
	config := openai.DefaultConfig(apiToken)
	config.BaseURL = target.ApiUrl

	client := openai.NewClientWithConfig(config)

	log.Printf("	🤖 Requesting chat completion at %s", process.Target)

	startTime := time.Now()
	resp, err := client.CreateChatCompletion(context.Background(), req)
	latency := time.Since(startTime).Milliseconds()

	log.Printf("	⏰ Took: %.2f seconds", float32(latency)/1000)

	if err != nil {
		return nil, fmt.Errorf("error requesting chat completion for %s: %v", process.Name, err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty response choices for %s", process.Name)
	}

	response := resp.Choices[0].Message.Content

	if len(response) == 0 {
		return nil, fmt.Errorf("missing response for %s: Empty response", process.Name)
	}

	if process.JSONSchema != nil {
		jsonObject, err := ExtractJsonFromText(response)
		if err != nil {
			// If error while extracting JSON try to parse the whole response
			jsonObj, errParsing := parseJSONWithSchema(response, process.JSONSchema)
			if errParsing != nil {
				return response, errParsing
			}

			return jsonObj, nil
		}

		// If JSON was found, then validate each JSON object
		var resultObjs = []any{}
		for _, jsonStr := range jsonObject {
			jsonObj, errParsing := parseJSONWithSchema(jsonStr, process.JSONSchema)
			if err != nil {
				return response, errParsing
			}

			resultObjs = append(resultObjs, jsonObj)
		}

		return resultObjs, nil
	}

	return response, nil
}

func parseJSONWithSchema(jsonStr string, schema JSONSchema) (any, error) {
	documentLoader := gojsonschema.NewStringLoader(jsonStr)
	schemaLoader := gojsonschema.NewGoLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		return nil, fmt.Errorf("error validating JSON schema")
	}

	if result.Valid() {
		var jsonObj any
		err := json.Unmarshal([]byte(jsonStr), &jsonObj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON")
		}

		return jsonObj, nil
	} else {
		return nil, fmt.Errorf("invalid JSON schema: %v", result.Errors())
	}
}

func parseTemplate(templateName, templateContent string, templateData interface{}) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("error parsing template %s: %v", templateName, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return "", fmt.Errorf("error executing template %s: %v", templateName, err)
	}

	result := buf.String()
	return result, nil
}
