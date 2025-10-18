package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/llm"
	"github.com/openai/openai-go/v3"
	"github.com/spf13/viper"
)

var itinerarySchema = openai.FunctionParameters{
	"$schema":  "http://json-schema.org/draft-07/schema#",
	"title":    "Itinerary",
	"type":     "object",
	"required": []string{"destination", "days"},
	"properties": map[string]any{
		"destination": map[string]any{"type": "string"},
		"startDate":   map[string]any{"type": "string"},
		"endDate":     map[string]any{"type": "string"},
		"currency":    map[string]any{"type": "string"},
		"days": map[string]any{
			"type": "array", "minItems": 1,
			"items": map[string]any{
				"type": "object", "required": []string{"day", "items"},
				"properties": map[string]any{
					"day":   map[string]any{"type": "integer", "minimum": 1},
					"label": map[string]any{"type": "string"},
					"items": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object", "required": []string{"title"},
							"properties": map[string]any{
								"title":     map[string]any{"type": "string"},
								"city":      map[string]any{"type": "string"},
								"place":     map[string]any{"type": "string"},
								"category":  map[string]any{"type": "string"},
								"startTime": map[string]any{"type": "string"},
								"endTime":   map[string]any{"type": "string"},
								"notes":     map[string]any{"type": "string"},
								"lat":       map[string]any{"type": "number"},
								"lon":       map[string]any{"type": "number"},
							},
							"additionalProperties": false,
						},
					},
				},
				"additionalProperties": false,
			},
		},
	},
	"additionalProperties": false,
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
		Client:  client,
		History: msgs,
		Model:   model,
	}
}

func (c *Controller) StreamMessage(
	ctx context.Context,
	userInput UserInput,
	chunkCh chan<- string,
) error {

	defer close(chunkCh) // close the channel at the end of this function

	// add user query to context history
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(userInput.Content),
			},
		},
	})

	tools := []openai.ChatCompletionToolUnionParam{
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        "save_itinerary",
					Description: openai.String("Call this ONLY when a finalized itinerary is ready."),
					Parameters: itinerarySchema,
				},
			},
		},
	}

	stream := c.Client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(c.Model.Name),
		Messages: c.History,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
		Tools: tools,
	})

	acc := openai.ChatCompletionAccumulator{}
	var tokenBuilder strings.Builder

	for stream.Next() { // returns false when stream ends
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			delta := chunk.Choices[0].Delta.Content
			tokenBuilder.WriteString(delta)
			chunkCh <- delta // add the response chunk to channel
		}

		// if _, ok := acc.JustFinishedToolCall(); ok {
		//     // Extract and invoke tool function.
		//     // Then feed the result back by calling the LLM again
		// }
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("chat streaming error: %w", err)
	}

	finalContent := tokenBuilder.String()

	toolCalls := acc.Choices[0].Message.ToolCalls

	if len(toolCalls) == 0 {
		fmt.Println("NO TOOL CALLED")
	}
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "save_itinerary" {
			fmt.Println("SAVE_ITINERARY TOOL CALLED!")
			var itin Itinerary
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &itin)
			if err != nil {
				panic(err)
			}

			c.Itinerary = &itin
		}
	}

	// Add assistant response to History
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			Content: openai.ChatCompletionAssistantMessageParamContentUnion{
				OfString: openai.String(finalContent),
			},
		},
	})

	return nil
}

func (c *Controller) GetHistory(chatID string) ([]openai.ChatCompletionMessageParamUnion, error) {
	// instead of keeping history as chatcompletion obj,
	// keep it as a []Message object
	// then write a function to convert []Message to chatcompletion
	// this should run once at the beginning of the chat session

	return c.History, nil
}


func (c *Controller) GetItinerary(chatID string) (*Itinerary, error) {
	// instead of keeping history as chatcompletion obj,
	// keep it as a []Message object
	// then write a function to convert []Message to chatcompletion
	// this should run once at the beginning of the chat session

	return c.Itinerary, nil
}


