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

	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

const (
	FunctionTypeName = "json_schema"
)

type JSONData struct {
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
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	for _, process := range config.Processes {
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

				log.Printf("💾 %s processing", debugInfo)

				data, err := requestModel(
					client,
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

				log.Printf("✅ %s processing completed!", debugInfo)
			}
			log.Printf("✅ Step(%d) completed!", i)
		}
		log.Printf("✅ All steps completed!")
	}

	log.Info("🎉🍻 All processes completed successfully!")

	return nil
}

func DumpDebugRawData(data JSONData, outputDir string) error {
	if data.Raw != "" {
		log.Errorf("💥 An error occurred, saving raw response to a file for debugging")
		error := DumpDataIntoOutputDir([]byte(data.Raw), outputDir, "debug")
		if error != nil {
			return error
		}
		log.Info("💾 Raw response saved")
	}
	return nil
}

func requestModel(client *openai.Client, process Process, corpus string, jsonSchema JSONSchema) (JSONData, error) {
	templateData := struct {
		Document string
	}{Document: string(corpus)}

	systemPrompt, err := executeTemplate("systemPrompt", process.SystemPrompt, templateData)
	if err != nil {
		return JSONData{}, fmt.Errorf("error executing system prompt template for %s: %v", process.Name, err)
	}

	userPrompt, err := executeTemplate("userPrompt", process.UserPrompt, templateData)
	if err != nil {
		return JSONData{}, fmt.Errorf("error executing user prompt template for %s: %v", process.Name, err)
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
		MaxTokens:   process.MaxTokens,
		Temperature: float32(process.Temperature),
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: openai.FunctionDefinition{
					Name:       FunctionTypeName,
					Parameters: jsonSchema,
				},
			},
		},
	}

	log.Printf("	🤖 Requesting chat completion")

	startTime := time.Now()
	resp, err := client.CreateChatCompletion(context.Background(), req)
	latency := time.Since(startTime).Milliseconds()

	log.Printf("	⏰ Took: %.2f seconds", float32(latency)/1000)

	if err != nil {
		return JSONData{}, fmt.Errorf("error requesting chat completion for %s: %v", process.Name, err)
	}

	var jsonString string

	if len(resp.Choices) == 0 {
		return JSONData{}, fmt.Errorf("empty response choices for %s", process.Name)
	}

	for _, toolCall := range resp.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == FunctionTypeName {
			jsonString = toolCall.Function.Arguments
			break
		}
	}

	if jsonString == "" {
		return JSONData{Raw: resp.Choices[0].Message.Content}, fmt.Errorf("missing `%s` function in response for %s", FunctionTypeName, process.Name)
	}

	response, err := parseJSONResponse(jsonString, process.Name)

	return JSONData{Raw: jsonString, Data: response}, err
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
