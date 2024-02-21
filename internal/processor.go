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
)

const (
	FunctionTypeName = "json_schema"
)

type ResponseData struct {
	Raw  string
	Data any
}

type Response struct {
	Conversations []Conversation `json:"conversations"`
}

type Conversation struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func Run(config Config, relativePath string) error {
	for _, process := range config.Processes {
		if process.Skip {
			log.Printf("🚫 Skipping: %s", process.Name)
			continue
		}

		outputDir := filepath.Join(relativePath, process.OutputDir)

		log.Printf("🤓 Processing: %s", process.Name)

		log.Printf("📚 Collecting corpus...")
		corpus, err := CollectCorpusData(process.Documents, relativePath)
		if err != nil {
			return err
		}

		log.Printf("📑 Corpus collected: %d", len(corpus))

		steps := 1
		if process.Steps != 0 {
			steps = process.Steps
		}

		for i := 1; i <= steps; i++ {
			stepInfo := fmt.Sprintf("🐾 Step(%d of %d)", i, steps)

			log.Printf("%s processing", stepInfo)
			chunkCorpus := SplitTextIntoChunks(corpus, process.ChunkSize)
			log.Printf("%s Chunking corpus into %d chunks", stepInfo, len(chunkCorpus))
			for chunkIndex, corpus := range chunkCorpus {
				chunkInfo := fmt.Sprintf("🍤 Chunk(%d of %d)", chunkIndex+1, len(chunkCorpus))
				debugInfo := fmt.Sprintf("%s %s ", stepInfo, chunkInfo)

				log.Printf("🥁 %s processing", debugInfo)

				target, err := getTargetByName(config.Targets, process.Target)
				if err != nil {
					return err
				}

				data, err := requestModel(
					target,
					process,
					string(corpus),
					process.JSONSchema,
				)

				if err != nil {
					// It might that the JSON is not formatted correctly, or an error was found
					// then save raw JSON to a file for debugging
					dumpError := DumpDebugRawData(data, outputDir)
					if dumpError != nil {
						return fmt.Errorf("error dumping raw data while: %v", err)
					}

					return err
				}

				if data.Data != nil {
					response, isOk := data.Data.(Response)
					if !isOk {
						err := DumpDebugRawData(data, outputDir)
						if err != nil {
							return err
						}
						return fmt.Errorf("error casting response might be because we got an empty JSON response %s", process.Name)
					}

					err = exportProcessResponseData(response, outputDir)
					if err != nil {
						return err
					}
				} else {
					err := DumpRawData(data.Raw, outputDir)
					if err != nil {
						return err
					}
				}

				log.Printf("	✅ %s processing completed!", debugInfo)
			}
			log.Printf("✅ Step(%d) completed!", i)
		}
		log.Printf("✅ All steps completed!")
	}

	log.Printf("🎉🍻 All processes completed successfully!")

	return nil
}

func DumpDebugRawData(data ResponseData, outputDir string) error {
	if data.Raw != "" {
		log.Print("💥 An error occurred, saving raw response to a file for debugging")
		error := DumpDataIntoOutputDir([]byte(data.Raw), outputDir, "debug")
		if error != nil {
			return error
		}
		log.Fatal("💾 Raw response saved for debugging")
	}
	return nil
}

func DumpRawData(rawData string, outputDir string) error {
	if len(rawData) == 0 {
		return nil
	}

	error := DumpDataIntoOutputDir([]byte(rawData), outputDir, "")
	if error != nil {
		return error
	}

	log.Printf("💾 Raw response saved")

	return nil
}

func getTargetByName(targets []Target, name string) (Target, error) {
	for _, target := range targets {
		if target.Name == name {
			return target, nil
		}
	}

	return Target{}, fmt.Errorf("target not found: %s", name)
}

func requestModel(target Target, process Process, corpus string, jsonSchema JSONSchema) (ResponseData, error) {
	templateData := struct {
		Document string
	}{Document: string(corpus)}

	systemPrompt, err := executeTemplate("systemPrompt", process.SystemPrompt, templateData)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error executing system prompt template for %s: %v", process.Name, err)
	}

	userPrompt, err := executeTemplate("userPrompt", process.UserPrompt, templateData)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error executing user prompt template for %s: %v", process.Name, err)
	}

	processJSONSchema := jsonSchema != nil && !process.SkipJsonSchema
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

	if processJSONSchema {
		req.Tools = []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: openai.FunctionDefinition{
					Name:       FunctionTypeName,
					Parameters: jsonSchema,
				},
			},
		}
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
		return ResponseData{}, fmt.Errorf("error requesting chat completion for %s: %v", process.Name, err)
	}

	if len(resp.Choices) == 0 {
		return ResponseData{}, fmt.Errorf("empty response choices for %s", process.Name)
	}

	var responseString string = ""

	for _, choice := range resp.Choices {
		if processJSONSchema {
			for _, toolCall := range choice.Message.ToolCalls {
				if toolCall.Function.Name == FunctionTypeName {
					responseString += toolCall.Function.Arguments
				}
			}
		} else {
			responseString += choice.Message.Content
		}
	}

	if processJSONSchema {
		if len(responseString) == 0 {
			log.Printf("	🙄 missing `%s` function in response for %s: Empty response.", FunctionTypeName, process.Name)
			return ResponseData{}, nil
		}

		response, err := parseJSONResponse(responseString, process.Name)
		return ResponseData{Raw: responseString, Data: response}, err
	}

	return ResponseData{Raw: responseString, Data: nil}, nil
}

func executeTemplate(templateName, templateContent string, templateData interface{}) (string, error) {
	tmpl := template.Must(template.New(templateName).Parse(templateContent))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, templateData)
	if err != nil {
		return "", fmt.Errorf("error executing template %s: %v", templateName, err)
	}

	result := buf.String()
	return result, nil
}

func parseJSONResponse(jsonString, promptName string) (Response, error) {
	var response Response
	rawJSONDataBytes := []byte(jsonString)
	err := json.Unmarshal(rawJSONDataBytes, &response)
	if err != nil {
		return Response{}, fmt.Errorf("error parsing JSON %s:%s", promptName, jsonString)
	}

	return response, nil
}

func getDocumentPaths(documents []string, relativePath string) ([]string, error) {
	uniqRelativeDocumentPaths := map[string]string{}
	for _, documentPath := range documents {
		relativeDocumentPath := filepath.Join(relativePath, documentPath)
		relativeDocumentPaths, err := ExpandFiles(relativeDocumentPath)
		if err != nil {
			return []string{}, err
		}

		for _, relativeDocumentPath := range relativeDocumentPaths {
			uniqRelativeDocumentPaths[relativeDocumentPath] = relativeDocumentPath
		}
	}

	var documentPaths []string
	for _, relativeDocumentPath := range uniqRelativeDocumentPaths {
		documentPaths = append(documentPaths, relativeDocumentPath)
	}

	return documentPaths, nil
}
