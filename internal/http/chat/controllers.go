package chat

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/llm"
	"github.com/openai/openai-go"
	"github.com/spf13/viper"
)

type Model struct {
	Name string `mapstructure:"model_name"`
	Type string `mapstructure:"model_type"`
	SystemPrompt string `mapstructure:"system_prompt"`
}

type Message struct {
	Role string
	Content string
}

type Controller struct {
	Client *openai.Client // OpenAI Client
	Model
	History []openai.ChatCompletionMessageParamUnion // context memory to store messages
}

func NewController() *Controller {
	if err := godotenv.Load(); err != nil {
		slog.Error("could not load model env " + err.Error())
	}

	client, err := llm.NewOpenAIClient()
	if err != nil {
		slog.Error(err.Error())
	}

	// loading model config
	configPath := os.Getenv("MODEL_CONFIG_PATH")

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := v.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	var model Model
	if err := v.Unmarshal(&model); err != nil {
		slog.Error("Failed to load model config file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if model.Name == "" {
		slog.Error("No model name specified")
	}

	// fmt.Print(model)

	if model.SystemPrompt == "" {
		slog.Error("No system prompt specified")
	}

	var msgs []openai.ChatCompletionMessageParamUnion
	msgs = append(msgs, openai.ChatCompletionMessageParamUnion{
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(model.SystemPrompt),
			},
		},
	})

	return &Controller{
		Client: client,
		History: msgs,
		Model: model,
	}
}

func (c *Controller) SendMessage(ctx context.Context, chatID, userID, content string) (*openai.ChatCompletion, error) {
	// add user query to context history
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(content),
			},
		},
	})


	resp, err := c.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(c.Model.Name),
		Messages: c.History,
	})

	// add assitant response to context history
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			Content: openai.ChatCompletionAssistantMessageParamContentUnion{
				OfString: openai.String(resp.Choices[0].Message.Content),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("chat api error: %w", err)
	}

	return resp, nil
}

func (c *Controller) GetHistory(chatID string) ([]openai.ChatCompletionMessageParamUnion, error) {
	return c.History, nil
}