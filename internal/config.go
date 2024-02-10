package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type JSONSchema map[string]interface{}

type Config struct {
	Processes []Process `yaml:"processes"`
}

type Process struct {
	Name         string   `yaml:"name"`
	Model        string   `yaml:"model"`
	Temperature  int      `yaml:"temperature"`
	MaxTokens    int      `yaml:"max_tokens"`
	ChunkSize    int      `yaml:"chunk_size"`
	Steps        int      `yaml:"steps"`
	OutputDir    string   `yaml:"output_dir"`
	Skip         bool     `yaml:"skip"`
	Documents    []string `yaml:"documents"`
	SystemPrompt string   `yaml:"system_prompt"`
	UserPrompt   string   `yaml:"user_prompt"`
	JSONSchema   `yaml:"json_schema"`
}

func NewConfig(path string) (Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	error := yaml.Unmarshal(data, &config)
	if error != nil {
		return Config{}, error
	}

	return config, nil
}
