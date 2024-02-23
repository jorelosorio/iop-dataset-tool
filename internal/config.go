package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type JSONSchema map[string]any

type Config struct {
	Processes []Process `yaml:"processes"`
	Targets   []Target  `yaml:"targets"`
}

type Target struct {
	Name      string `yaml:"name"`
	ApiUrl    string `yaml:"api_url"`
	ApiKeyEnv string `yaml:"api_key_env"`
}

type Process struct {
	Name         string   `yaml:"name"`
	Model        string   `yaml:"model"`
	Target       string   `yaml:"target"`
	Temperature  float32  `yaml:"temperature"`
	MaxTokens    int      `yaml:"max_tokens"`
	ChunkSize    int      `yaml:"chunk_size"`
	OutputDir    string   `yaml:"output_dir"`
	Skip         bool     `yaml:"skip"`
	Documents    []string `yaml:"documents"`
	SystemPrompt string   `yaml:"system_prompt"`
	UserPrompt   string   `yaml:"user_prompt"`
	JSONSchema   `yaml:"json_schema,inline"`
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

	newConfig, err := applyDefaultValues(config)
	if err != nil {
		return Config{}, error
	}

	return newConfig, nil
}

func (c Config) GetTarget(name string) (Target, error) {
	for _, target := range c.Targets {
		if target.Name == name {
			return target, nil
		}
	}

	return Target{}, fmt.Errorf("target not found: %s", name)
}

func applyDefaultValues(config Config) (Config, error) {
	var emptyConfig Config
	// Apply default values to Targets
	for i, target := range config.Targets {
		if target.Name == "" {
			return emptyConfig, fmt.Errorf("name is required for target: %d", i)
		}
		if target.ApiUrl == "" {
			return emptyConfig, fmt.Errorf("api_url is required for target: %s", target.Name)
		}
		if target.ApiKeyEnv == "" {
			return emptyConfig, fmt.Errorf("api_key_env is required for target: %d", i)
		}

		if os.Getenv(target.ApiKeyEnv) == "" {
			return emptyConfig, fmt.Errorf("api_key_env %s is not set on .env or system env", target.ApiKeyEnv)
		}
	}

	var process *Process
	for i := range config.Processes {
		process = &config.Processes[i]

		if process.Name == "" {
			return emptyConfig, fmt.Errorf("name is required for process: %d", i)
		}

		if process.Model == "" {
			return emptyConfig, fmt.Errorf("model is required for process: %s", process.Name)
		}

		if process.Target == "" {
			return emptyConfig, fmt.Errorf("target is required for process: %s", process.Name)
		}

		_, err := config.GetTarget(process.Target)
		if err != nil {
			return emptyConfig, err
		}

		if process.MaxTokens == 0 {
			process.MaxTokens = 4096
		}

		if process.ChunkSize == 0 {
			process.ChunkSize = 2048
		}

		if process.OutputDir == "" {
			process.OutputDir = "output"
		}

		if process.SystemPrompt == "" {
			return emptyConfig, fmt.Errorf("system_prompt is required for process: %s", process.Name)
		}

		if process.UserPrompt == "" {
			return emptyConfig, fmt.Errorf("user_prompt is required for process: %s", process.Name)
		}
	}

	return config, nil
}
