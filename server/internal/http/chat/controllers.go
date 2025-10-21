package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	calltools "github.com/nakul-krishnakumar/kaiyo-ai/internal/http/chat/call_tools"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/llm"
	"github.com/openai/openai-go/v3"
	"github.com/spf13/viper"
)

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

	var tools = calltools.NewCallTools()

	return &Controller{
		Client:  client,
		History: msgs,
		Model:   model,
		Tools:   tools,
	}
}

func (c *Controller) doToolCalling(ctx context.Context) error {
	params := openai.ChatCompletionNewParams{
		Messages:        c.History,
		Tools:           c.Tools.InitOpenAITools(),
		Seed:            openai.Int(0),
		Model:           openai.ChatModel(c.Model.Name),
		ReasoningEffort: openai.ReasoningEffortMedium,
	}

	completion, err := c.Client.Chat.Completions.New(ctx, params)
	if err != nil {
		return err
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	// if no tools called, then stop the chat
	if len(toolCalls) == 0 {
		fmt.Println("NO TOOL CALLED")
		c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(completion.Choices[0].Message.Content),
				},
			},
		})

		return nil
	}

	c.History = append(c.History, completion.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		result := c.Tools.HandleToolCall(ctx, toolCall.Function.Name, toolCall.Function.Arguments)

		// Append tool result to history
		resultJSON, _ := json.Marshal(result)
		c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
			OfTool: &openai.ChatCompletionToolMessageParam{
				ToolCallID: toolCall.ID,
				Content: openai.ChatCompletionToolMessageParamContentUnion{
					OfString: openai.String(string(resultJSON)),
				},
			},
		})
	}

	return nil
}

func (c *Controller) performPlanningPhase(ctx context.Context, userInput UserInput) error {
	// Add user query to history
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(userInput.Content),
			},
		},
	})

	const maxIters = 3
	for iter := 0; iter < maxIters; iter++ {
		if err := c.doToolCalling(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) performStreamingPhase(ctx context.Context, chunkCh chan<- string) error {
	defer close(chunkCh)
	fmt.Println("INSIDE PHASE 2") /////////////////////////////////////

	// Add a user prompt to trigger narrative generation
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(`Generate the complete itinerary in proper Markdown format, ONLY if all necessary data is available.

CRITICAL FORMATTING RULES:
1. Add TWO newlines (\n\n) after every heading (##)
2. Add TWO newlines (\n\n) before and after horizontal rules (---)
3. Add ONE newline (\n) after each bullet point
4. Add TWO newlines (\n\n) between different sections
5. Use proper markdown syntax:
   - Headings: ## Day 1: Title
   - Horizontal rules: ---
   - Bullet points: - Item text
   - Bold: **text**
   - Italic: *text*

Example format:
## Day 1: Monday, April 22 – Welcome to Paris

- 09:00 – 11:00 • Activity Name
  Description here
- 11:00 – 14:00 • Next Activity
  Description here

---

## Day 2: Tuesday, April 23 – Next Day

- 08:30 – 10:00 • Morning Activity
  Description here

Generate the itinerary now with proper spacing.`),
			},
		},
	})

	stream := c.Client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(c.Model.Name),
		Messages: c.History,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
		// Tools: tools,
	})

	acc := openai.ChatCompletionAccumulator{}
	var tokenBuilder strings.Builder

	for stream.Next() { // returns false when stream ends
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)

			if chunk.Choices[0].Delta.Content != "" {
				delta := chunk.Choices[0].Delta.Content
				tokenBuilder.WriteString(delta)
				chunkCh <- delta // add the response chunk to channel
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("chat streaming error: %w", err)
	}

	finalContent := tokenBuilder.String()

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

func (c *Controller) performSavingPhase() error {
	// Add a user prompt to trigger save_itinerary tool call
	c.History = append(c.History, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(`
				IF AND ONLY IF all the necessary details for an itinerary exists, call save_itinerary tool to save the data.`),
			},
		},
	})

	tools := []openai.ChatCompletionToolUnionParam{
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        "save_itinerary",
					Description: openai.String("Call this ONLY after narrative, to save the finalised itinerary."),
					Parameters:  itinerarySchema,
				},
			},
		},
	}

	params := openai.ChatCompletionNewParams{
		Messages:        c.History,
		Tools:           tools,
		Seed:            openai.Int(0),
		Model:           openai.ChatModel(c.Model.Name),
		ReasoningEffort: openai.ReasoningEffortMedium,
	}

	completion, err := c.Client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		return err
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

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

	return nil
}

func (c *Controller) StreamMessage(ctx context.Context, userInput UserInput, chunkCh chan<- string) error {

	// Synchronous tool orchestration
	if err := c.performPlanningPhase(ctx, userInput); err != nil {
		return err
	}

	fmt.Println("PHASE 1 DONE")

	// Streaming final response
	if err := c.performStreamingPhase(ctx, chunkCh); err != nil {
		return err
	}

	// save_itinerary tool is called as a background job
	if err := c.performSavingPhase(); err != nil {
		slog.Error("Could not save itinerary", slog.String("error", err.Error()))
	}

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
