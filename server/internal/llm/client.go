package llm

import (
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
)

// Config holds Azure OpenAI configuration
type Config struct {
	apiKey     string
	apiVersion string
	endpoint   string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	apiKey := os.Getenv("AZURE_OPEN_AI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("AZURE_OPEN_AI_API_KEY environment variable is required")
	}

	endpoint := os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("AZURE_OPEN_AI_ENDPOINT environment variable is required")
	}

	apiVersion := os.Getenv("AZURE_OPEN_AI_API_VERSION")
	if apiVersion == "" {
		apiVersion = "2024-08-01-preview" // default value
	}

	return &Config{
		apiKey:     apiKey,
		apiVersion: apiVersion,
		endpoint:   endpoint,
	}, nil
}

// NewClient creates a new Azure OpenAI client
func NewOpenAIClient() (*openai.Client, error) {
	config, err := LoadFromEnv()

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	client := openai.NewClient(
		option.WithAPIKey(config.apiKey),
		azure.WithEndpoint(config.endpoint, config.apiVersion),
	)

	return &client, nil
}
